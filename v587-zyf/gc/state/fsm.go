package state

import (
	"errors"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"time"
)

type Fsm struct {
	curState    iface.IState
	transitions map[iface.StateTrigger]map[iface.IState]iface.IState
	isRun       bool
	sleepTime   int64
}

func NewFsm() *Fsm {
	return &Fsm{
		transitions: make(map[iface.StateTrigger]map[iface.IState]iface.IState),
		isRun:       false,
	}
}

func (f *Fsm) SetSleepTime(t int) {
	f.sleepTime = time.Now().UnixMilli() + int64(t)
}

func (f *Fsm) Start(initState iface.IState) {
	f.curState = initState
	initState.Enter()
	initState.Execute()
}

func (f *Fsm) Run() {
	if f.curState != nil && f.isRun {
		if time.Now().UnixMilli() > f.sleepTime {
			f.curState.ExecuteBefore()
			f.curState.Execute()
		}
	}
}

func (f *Fsm) Stop() {
	if f.curState != nil {
		f.curState.End()
	}
}

func (f *Fsm) Pause() {
	f.isRun = false
}

func (f *Fsm) Recover() {
	f.isRun = true
}

func (f *Fsm) GetRunState() bool {
	return f.isRun
}

func (f *Fsm) Register(trigger iface.StateTrigger, fromStates []iface.IState, toState iface.IState) {
	var stateMap map[iface.IState]iface.IState
	var ok bool
	for _, fromState := range fromStates {
		if stateMap, ok = f.transitions[trigger]; !ok {
			stateMap = make(map[iface.IState]iface.IState)
			f.transitions[trigger] = stateMap
		}
		stateMap[fromState] = toState
	}
}

func (f *Fsm) SetInitState(state iface.IState) {
	f.curState = state
}

func (f *Fsm) Event(trigger iface.StateTrigger) error {
	if stateMap, ok := f.transitions[trigger]; ok {
		if state, ok := stateMap[f.curState]; ok {
			f.curState.End()
			f.curState = state
			state.Enter()
			return nil
		}
	}
	return errors.New("no state switch found")
}

func (f *Fsm) Debug() {
	log.Debug("fsm", zap.Bool("isRun", f.isRun), zap.Bool("curState", f.curState != nil), zap.Bool("time", time.Now().UnixMilli() > f.sleepTime))
}
