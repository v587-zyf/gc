package event

import (
	"context"
	"github.com/v587-zyf/gc/log"
	"testing"

	"go.uber.org/zap"
)

func init() {
	log.Init(context.Background())
}

func TestEventEmitter(t *testing.T) {
	p := NewPool()

	e, _ := p.Get(123)

	fn1 := func(a int, b string, c int32) {
		log.Debug("callback 1", zap.Int("a", a), zap.String("b", b), zap.Int32("c", c))
	}

	fn2 := func(a int, b string, c int32) {
		log.Debug("callback 2")
	}

	wrapf1, _ := GenListener(fn1)
	wrapf2, _ := GenListener(fn2)
	e.On("abc", wrapf1)
	e.Once("aa", wrapf2)

	e.Emit("abc", 10, "abc", int32(123))
}
