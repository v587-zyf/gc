package iface

import (
	"context"
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type ITcpSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() net.Conn
	GetCtx() context.Context

	Send(msgID uint16, tag uint32, userID uint64, msg IProtoMessage) error
}

type IWsSession interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Remove(key string)

	GetID() uint64
	SetID(id uint64)

	Start()
	Close() error

	GetConn() *websocket.Conn
	GetCtx() context.Context

	SendMsg(fn func(args ...any) ([]byte, error), args ...any) error

	Login()
	Heartbeat()
	IsHeartbeatTimeout(time time.Time) bool
}

type ITcpSessionMethod interface {
	Start(ss ITcpSession)
	Recv(conn ITcpSession, data any)
	Stop(ss ITcpSession)
}
type IWsSessionMethod interface {
	Start(ss IWsSession)
	Recv(conn IWsSession, data any)
	Stop(ss IWsSession)
}
