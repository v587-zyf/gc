package ws_session

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"time"
)

type Session struct {
	id   uint64
	conn *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc

	hooks  *Hooks
	cache  sync.Map
	method iface.IWsSessionMethod

	outChan chan []byte

	once sync.Once

	heartbeatTime time.Time
}

func NewSession(ctx context.Context, conn *websocket.Conn) *Session {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ctx:    ctx,
		cancel: cancel,

		outChan: make(chan []byte, 1024),

		hooks: NewHooks(),

		heartbeatTime: time.Now(),
	}
	s.conn = conn

	return s
}

func (s *Session) Start() {
	s.hooks.ExecuteStart(s)

	go s.readPump()
	go s.IOPump()
}

func (s *Session) Hooks() *Hooks {
	return s.hooks
}

func (s *Session) Set(key string, value any) {
	s.cache.Store(key, value)
}
func (s *Session) Get(key string) (any, bool) {
	return s.cache.Load(key)
}
func (s *Session) Remove(key string) {
	s.cache.Delete(key)
}

func (s *Session) GetID() uint64 {
	return s.id
}
func (s *Session) SetID(id uint64) {
	if id <= 0 {
		id = 0
	}
	s.id = id
}

func (s *Session) Close() error {
	s.once.Do(func() {
		s.hooks.ExecuteStop(s)

		close(s.outChan)

		s.cancel()
		s.conn.Close()

		sessionMgr.unRegisterCh <- s
	})

	return nil
}

func (s *Session) GetConn() *websocket.Conn {
	return s.conn
}

func (s *Session) GetCtx() context.Context {
	return s.ctx
}

func (s *Session) Login() {
	sessionMgr.loginCh <- s
}

func (s *Session) Heartbeat() {
	s.heartbeatTime = time.Now()
}

func (s *Session) IsHeartbeatTimeout(time time.Time) bool {
	return s.heartbeatTime.Add(enums.HEARTBEAT_TIMEOUT).Before(time)
}

func (s *Session) SendMsg(fn func(args ...any) ([]byte, error), args ...any) error {
	sendBytes, err := fn(args...)
	if err != nil {
		return err
	}

	select {
	case s.outChan <- sendBytes:
		return nil
	default:
		return errcode.ERR_NET_SEND_TIMEOUT
	}
}

func (s *Session) readPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("readMsgLoop panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
		}
	}()

LOOP:
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				//log.Info("ws_server read err", zap.Error(err))
			}
			break LOOP
		}
		if len(message) > enums.MAX_MSG_SIZE {
			log.Warn("消息超过最大长度，忽略", zap.Int("length", len(message)))
			continue
		}

		if message != nil && len(message) > 0 {
			dataCopy := make([]byte, len(message))
			copy(dataCopy, message)
			s.hooks.ExecuteRecv(s, dataCopy)
		}
	}

	s.cancel()
}

func (s *Session) IOPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("IOPump panic", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
		}
	}()

LOOP:
	for {
		select {
		case data := <-s.outChan:
			s.conn.SetWriteDeadline(time.Now().Add(enums.CONN_WRITE_WAIT_TIME))

			if err := s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				//msgID := binary.BigEndian.Uint16(data[0:2])
				//log.Warn("conn write err", zap.Uint64("userID", s.id),
				//	zap.Uint16("msgID", msgID), zap.Int("len", len(data)), zap.Error(err))
				break LOOP
			}
		case <-s.ctx.Done():
			break LOOP
		}
	}

	s.Close()
}
