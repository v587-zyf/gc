package ws_server

import (
	"github.com/v587-zyf/gc/iface"
)

type WsOption struct {
	addr string
	pem  string
	key  string

	https bool

	method iface.IWsSessionMethod
}

type Option func(opts *WsOption)

func NewWsOption() *WsOption {
	o := &WsOption{}

	return o
}

func WithAddr(addr string) Option {
	return func(opts *WsOption) {
		opts.addr = addr
	}
}

func WithPem(pem string) Option {
	return func(opts *WsOption) {
		opts.pem = pem
	}
}

func WithKey(key string) Option {
	return func(opts *WsOption) {
		opts.key = key
	}
}

func WithMethod(m iface.IWsSessionMethod) Option {
	return func(opts *WsOption) {
		opts.method = m
	}
}

func WithHttps(https bool) Option {
	return func(opts *WsOption) {
		opts.https = https
	}
}
