package iface

import (
	"context"
	"github.com/gorilla/websocket"
)

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

	Send2User(msgID int32, msg IProtoMessage) error
	SendErr2User(err error) error
}

type IWsSessionMethod interface {
	Start(ss IWsSession)
	Recv(conn IWsSession, data any)
	Stop(ss IWsSession)
}
