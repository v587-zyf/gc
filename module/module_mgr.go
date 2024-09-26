package module

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"sync"
)

type ModuleMgr struct {
	modules sync.Map // name:iface.IModule
}

func (mm *ModuleMgr) Add(m iface.IModule) {
	if _, ok := mm.modules.Load(m.Name()); ok {
		log.Warn("module already exists", zap.String("name", m.Name()))
		return
	}

	mm.modules.Store(m.Name(), m)
}

func (mm *ModuleMgr) Get(name string) iface.IModule {
	if module, ok := mm.modules.Load(name); ok {
		return module.(iface.IModule)
	}

	log.Warn("module not exists", zap.String("name", name))
	return nil
}

func (mm *ModuleMgr) Del(name string) {
	if _, ok := mm.modules.Load(name); ok {
		mm.modules.Delete(name)
		log.Info("module deleted", zap.String("name", name))
	} else {
		log.Warn("module not found for deletion", zap.String("name", name))
	}
}

func (mm *ModuleMgr) Length() int {
	length := 0
	mm.modules.Range(func(key, value any) bool {
		length++
		return true
	})

	return length
}

func (mm *ModuleMgr) Init(ctx context.Context, opts ...iface.Option) (err error) {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return fmt.Errorf("no module")
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go func(module iface.IModule) {
			defer wg.Done()
			err := module.Init(ctx, opts...)
			if err != nil {
				log.Error("module init failed", zap.String("name", module.Name()), zap.Error(err))
			}
		}(value.(iface.IModule))
		return true
	})
	wg.Wait()

	return nil
}

func (mm *ModuleMgr) Start() (err error) {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return fmt.Errorf("no module")
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go func(module iface.IModule) {
			defer wg.Done()
			err := module.Start()
			if err != nil {
				log.Error("module start failed", zap.String("name", module.Name()), zap.Error(err))
			}
		}(value.(iface.IModule))
		return true
	})

	wg.Wait()

	return nil
}

func (mm *ModuleMgr) Run() {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go func(module iface.IModule) {
			defer wg.Done()
			module.Run()
		}(value.(iface.IModule))
		return true
	})

	wg.Wait()
}

func (mm *ModuleMgr) Stop() {
	moduleLen := mm.Length()
	if moduleLen <= 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(moduleLen)

	mm.modules.Range(func(key, value any) bool {
		go func(module iface.IModule) {
			defer wg.Done()
			module.Stop()
		}(value.(iface.IModule))
		return true
	})

	wg.Wait()
}
