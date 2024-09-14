package workerpool

import (
	"context"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

type WorkerPool struct {
	ctx     context.Context
	cancel  context.CancelFunc
	options *WorkerPoolOption

	maxWorkerCnt int
	minWorkerCnt int
	curWorkerCnt int

	mu       sync.Mutex
	pool     sync.Pool
	stopCh   chan struct{}
	mustStop bool

	ready         []*worker
	idleCleanTime time.Duration
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		maxWorkerCnt:  256 * 1024,
		minWorkerCnt:  10,
		idleCleanTime: 5 * 60 * time.Second,
		options:       NewWorkerPoolOption(),

		pool: sync.Pool{
			New: func() interface{} {
				return &worker{
					task: make(chan iface.ITask, workerCap),
				}
			},
		},

		ready: make([]*worker, 0),
	}
}

func (p *WorkerPool) Init(ctx context.Context, opts ...any) error {
	p.ctx, p.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(p.options)
		}
	}

	if p.options.maxCount != 0 {
		p.maxWorkerCnt = p.options.maxCount
	}

	return nil
}

var workerCap = func() int {
	// Use blocking workerChan if GOMAXPROCS=1.
	// This immediately switches Serve to WorkerFunc, which results
	// in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking workerChan if GOMAXPROCS>1,
	// since otherwise the Serve caller (Acceptor) may lag accepting
	// new connections if WorkerFunc is CPU-bound.
	return 1
}()

func (p *WorkerPool) Start() {
	if p.stopCh != nil {
		return
	}

	p.stopCh = make(chan struct{})
	stopCh := p.stopCh

	p.mustStop = false

	go func() {
		timer := time.NewTicker(10 * time.Second)
		var scratch []*worker
	LOOP:
		for {
			select {
			case <-stopCh:
				break LOOP
			case <-timer.C:
				p.clean(&scratch)
			}
		}

		timer.Stop()
		timer = nil
	}()
}

func (p *WorkerPool) Stop() {
	if p.stopCh == nil {
		return
	}
	close(p.stopCh)
	p.stopCh = nil

	p.mu.Lock()
	ready := p.ready
	for i := range ready {
		ready[i].task <- nil
		// close(ready[i].task)
		ready[i] = nil
	}
	p.ready = ready[:0]
	p.mustStop = true
	p.mu.Unlock()
}

func (p *WorkerPool) GetCtx() context.Context {
	return p.ctx
}

func (p *WorkerPool) Assign(task iface.ITask) error {
	w := p.getWorker()
	if w == nil {
		log.Error("WorkPool Assign Failed", zap.Int("maxWorkerCnt", p.maxWorkerCnt), zap.Any("task", task))
		return errcode.ERR_WP_TOO_MANY_WORKER
	}
	w.task <- task
	return nil
}

func (p *WorkerPool) clean(scratch *[]*worker) {
	criticalTime := time.Now().Add(-p.idleCleanTime)

	p.mu.Lock()
	ready := p.ready
	n := len(ready)
	if n <= p.minWorkerCnt {
		p.mu.Unlock()
		return
	}

	// Use binary-search algorithm to find out the index of the least recently worker which can be cleaned up.
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = (l + r) / 2
		if criticalTime.After(p.ready[mid].lastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	i := r
	if i == -1 {
		p.mu.Unlock()
		return
	}

	*scratch = append((*scratch)[:0], ready[:i+1]...)
	m := copy(ready, ready[i+1:])
	for i = m; i < n; i++ {
		ready[i] = nil
	}
	p.ready = ready[:m]
	p.mu.Unlock()

	// Notify obsolete workers to stop.
	// This notification must be outside the wp.lock, since ch.ch
	// may be blocking and may consume a lot of time if many workers
	// are located on non-local CPUs.
	tmp := *scratch
	for i := range tmp {
		tmp[i].task <- nil
		tmp[i] = nil
	}
	log.Info("tmp clean", zap.Int("len", len(tmp)))
}

func (p *WorkerPool) getWorker() *worker {
	var w *worker
	canCreate := false

	p.mu.Lock()
	ready := p.ready
	n := len(ready) - 1
	if n < 0 {
		if p.curWorkerCnt < p.maxWorkerCnt {
			canCreate = true
			p.curWorkerCnt++
		}
	} else {
		w = ready[n]
		ready[n] = nil
		p.ready = ready[:n]
	}
	p.mu.Unlock()

	if w == nil {
		if !canCreate {
			return nil
		}
		v := p.pool.Get()
		w = v.(*worker)
		go func() {
			w.run(p)

			p.mu.Lock()
			p.curWorkerCnt--
			p.mu.Unlock()

			p.pool.Put(v)
		}()
	}
	return w
}

func (p *WorkerPool) release(w *worker) bool {
	w.lastUseTime = time.Now()
	p.mu.Lock()
	if p.mustStop {
		p.mu.Unlock()
		return false
	}
	p.ready = append(p.ready, w)
	p.mu.Unlock()
	return true
}
