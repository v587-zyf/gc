package ws_server

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"net/http"
)

type WsServer struct {
	options *WsOption

	ctx    context.Context
	cancel context.CancelFunc

	upGrader *websocket.Upgrader
}

func NewWsServer() *WsServer {
	s := &WsServer{
		options: NewWsOption(),
	}

	return s
}

func (s *WsServer) Init(ctx context.Context, option ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(s.options)
	}

	s.upGrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return nil
}

func (s *WsServer) Start() {
	http.HandleFunc("/ws", s.wsHandle)
	var err error
	if s.options.dev {
		err = http.ListenAndServe(s.options.addr, nil)
	} else {
		err = http.ListenAndServeTLS(s.options.addr, s.options.pem, s.options.key, nil)
	}
	if err != nil {
		panic(err)
	}
}

func (s *WsServer) wsHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	wsConn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("webSocket upgrade err:", zap.Error(err))
	}

	ss := ws_session.NewSession(context.Background(), wsConn)
	ss.Hooks().OnMethod(s.options.method)
	ss.Start()
}

func (s *WsServer) Stop() {

}
func (s *WsServer) Wait() {

}
