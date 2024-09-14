package workerpool

import (
	"context"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
)

var defaultWorkPoll *WorkerPool

func Init(ctx context.Context, opts ...any) (err error) {
	defaultWorkPoll = NewWorkerPool()
	if err = defaultWorkPoll.Init(ctx, opts...); err != nil {
		return err
	}
	defaultWorkPoll.Start()

	return
}
func GetCtx() context.Context {
	return defaultWorkPoll.GetCtx()
}

func Assign(task iface.ITask) error {
	return defaultWorkPoll.Assign(task)
}

func AssignWsTask(fn ws_session.Recv, ss iface.IWsSession, data any) error {
	return defaultWorkPoll.Assign(&WsTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}
