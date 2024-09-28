package buffer_pool

import (
	"context"
)

var defaultBufferPool *BufferPool

func Init(ctx context.Context, opts ...any) (err error) {
	defaultBufferPool = NewBufferPool()
	if err = defaultBufferPool.Init(ctx, opts...); err != nil {
		return err
	}

	return
}
func GetCtx() context.Context {
	return defaultBufferPool.GetCtx()
}
func GetBufferPool() *BufferPool {
	return defaultBufferPool
}
func GetBuffer() *Buffer {
	return defaultBufferPool.Get()
}
func Put(buf *Buffer) {
	defaultBufferPool.Put(buf)
}

func GetStats() map[string]int {
	return defaultBufferPool.GetStats()
}
