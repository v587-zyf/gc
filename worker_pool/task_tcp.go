package worker_pool

import (
	"github.com/v587-zyf/gc/gcnet/tpc_session"
	"github.com/v587-zyf/gc/iface"
)

func (p *WorkerPool) AssignNetTask(fn tpc_session.Recv, ss iface.ITcpSession, data any) error {
	return Assign(&TcpTask{
		Func:    fn,
		Session: ss,
		Data:    data,
	})
}

type TcpTask struct {
	Func    tpc_session.Recv
	Session iface.ITcpSession
	Data    any
}

func (t *TcpTask) Do() {
	t.Func(t.Session, t.Data)
}
