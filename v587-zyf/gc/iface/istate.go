package iface

type IState interface {
	Enter()
	ExecuteBefore()
	Execute()
	ExecuteAfter()
	End()
}

type StateTrigger int
