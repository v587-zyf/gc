package ws_server

import (
	"context"
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/telegram/go_tg_bot"
	"go.uber.org/zap"
	"io"
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

func (s *WsServer) Init(ctx context.Context, option ...Option) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt(s.options)
	}

	s.upGrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return nil
}

func (s *WsServer) Start() {
	go ws_session.GetSessionMgr().Start()

	var err error

	if s.options.handler == nil {
		r := mux.NewRouter()
		r.HandleFunc("/api/test", s.test).Methods("GET")
		r.HandleFunc("/api/webHook", s.webHook).Methods("POST")
		r.HandleFunc("/ws", s.wsHandle).Methods("GET")
		s.options.handler = r
	}
	if len(s.options.handlerFuncs) > 0 {
		for _, v := range s.options.handlerFuncs {
			s.options.handler.(*mux.Router).HandleFunc(v.path, v.fn).Methods(v.methods)
		}
	}

	if s.options.https {
		err = http.ListenAndServeTLS(s.options.addr, s.options.pem, s.options.key, s.options.handler)
	} else {
		err = http.ListenAndServe(s.options.addr, s.options.handler)
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
	ws_session.GetSessionMgr().RegisterCh <- ss
	ss.Start()
}

func (s *WsServer) webHook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var update gotgbot.Update
	if err = json.Unmarshal(body, &update); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err = go_tg_bot.ProcessUpdate(&update); err != nil {
		log.Error("tg bot process update err:", zap.Error(err))
		http.Error(w, "Error process update", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *WsServer) test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *WsServer) Stop() {

}
