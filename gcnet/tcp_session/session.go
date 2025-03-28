package tcp_session

import (
	"context"
	"encoding/binary"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Session struct {
	id   uint64
	conn net.Conn

	ctx    context.Context
	cancel context.CancelFunc

	hooks  *Hooks
	cache  sync.Map
	method iface.ITcpSessionMethod

	inChan  chan []byte
	outChan chan []byte

	once sync.Once
}

func NewSession(ctx context.Context, conn net.Conn) *Session {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ctx:    ctx,
		cancel: cancel,

		inChan:  make(chan []byte, 1024),
		outChan: make(chan []byte, 1024),

		hooks: NewHooks(),
	}
	s.conn = conn

	return s
}

func (s *Session) Start() {
	go func() {
		s.readPump()
	}()

	go func() {
		s.parsePump()
	}()

	go func() {
		s.writePump()
	}()
}

func (s *Session) Hooks() *Hooks {
	return s.hooks
}

func (s *Session) Set(key string, value any) {
	s.cache.Store(key, value)
}
func (s *Session) Get(key string) (any, bool) {
	v, ok := s.cache.Load(key)
	if !ok {
		v = nil
	}
	return v, ok
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
		s.cancel()
		s.conn.Close()
	})

	return nil
}

func (s *Session) GetConn() net.Conn {
	return s.conn
}

func (s *Session) GetCtx() context.Context {
	return s.ctx
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
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			if err, ok := r.(error); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err.Error()), zap.ByteString("core", buf))
			} else if err, ok := r.(string); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err), zap.ByteString("core", buf))
			} else {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.Reflect("err", err), zap.ByteString("core", buf))
			}
		}
	}()

	s.hooks.ExecuteStart(s)

	//	scanner := bufio.NewScanner(s.conn)
	//	scanner.Buffer(make([]byte, enums.READ_BUFF_SIZE_INIT), enums.READ_BUFF_SIZE_MAX)
	//	scanner.Split(s.split)
	//LOOP:
	//	for {
	//		ok := scanner.Scan()
	//		if !ok {
	//			if err := scanner.Err(); err != nil {
	//				log.Error("server read err", zap.Error(err))
	//			}
	//			break LOOP
	//		}
	//
	//		data := scanner.Bytes()
	//		if len(data) > 0 {
	//			dataCopy := make([]byte, len(data))
	//			copy(dataCopy, data)
	//			select {
	//			case s.inChan <- dataCopy:
	//			default:
	//				log.Warn("inChan is full, dropping message", zap.Uint64("sessID", s.GetID()))
	//			}
	//		}
	//	}

	buffer := make([]byte, enums.READ_BUFF_SIZE_INIT)
LOOP:
	for {
		n, err := s.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			//log.Error("server read err", zap.Error(err))
			break LOOP
		}

		data := buffer[:n]
		if len(data) > 0 {
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)
			select {
			case s.inChan <- dataCopy:
			default:
				log.Warn("inChan is full, dropping message", zap.Uint64("sessID", s.GetID()))
			}
		}
	}

	s.cancel()
}

func (s *Session) parsePump() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			if err, ok := r.(error); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err.Error()), zap.ByteString("core", buf))
			} else if err, ok := r.(string); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err), zap.ByteString("core", buf))
			} else {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.Reflect("err", err), zap.ByteString("core", buf))
			}
		}
	}()

LOOP:
	for {
		select {
		case data := <-s.inChan:
			go func(d []byte) {
				s.hooks.ExecuteRecv(s, d)
			}(data)
		case <-s.ctx.Done():
			break LOOP
		}
	}
}

func (s *Session) writePump() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			if err, ok := r.(error); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err.Error()), zap.ByteString("core", buf))
			} else if err, ok := r.(string); ok {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.String("err", err), zap.ByteString("core", buf))
			} else {
				log.Error("core dump", zap.Uint64("sessID", s.GetID()),
					zap.Reflect("err", err), zap.ByteString("core", buf))
			}
		}
	}()

LOOP:
	for {
		select {
		case data := <-s.outChan:
			s.conn.SetWriteDeadline(time.Now().Add(enums.CONN_WRITE_WAIT_TIME))

			_, err := s.conn.Write(data)
			if err != nil {
				msgID := binary.BigEndian.Uint16(data[0:2])
				log.Warn("conn write err", zap.Uint64("userID", s.id),
					zap.Uint16("msgID", msgID), zap.Int("len", len(data)), zap.Error(err))
				break LOOP
			}
		case <-s.ctx.Done():
			break LOOP
		}
	}

	s.conn.Close()
	s.cancel()

	s.hooks.ExecuteStop(s)
}

func (s *Session) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	dataLen := len(data)
	if dataLen < enums.MSG_HEADER_SIZE {
		return 0, nil, nil
	}

	// body len
	n := int(binary.LittleEndian.Uint32(data[0:4]))
	//if n > enums.MSG_MAX_PACKET_SIZE-enums.MSG_HEADER_SIZE || n < 0 {
	//	log.Error("body len invalid", zap.Uint64("sessID", s.id),
	//		zap.Int("n", n), zap.String("addr", s.GetConn().RemoteAddr().String()))
	//	return 0, nil, errors.New("body len invalid")
	//}
	if dataLen < n+enums.MSG_HEADER_SIZE {
		return 0, nil, nil
	}

	return n + enums.MSG_HEADER_SIZE, data[:n+enums.MSG_HEADER_SIZE], nil
}
