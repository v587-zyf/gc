package ws_session_mgr

import (
	"github.com/v587-zyf/gc/iface"
	"sync"
	"time"
)

var sessionMgr *SessionMgr

func init() {
	sessionMgr = NewSessionMgr()
}

type SessionMgr struct {
	clients     map[iface.IWsSession]struct{}
	clientsLock sync.RWMutex

	online     map[uint64]iface.IWsSession
	onlineLock sync.RWMutex

	RegisterCh   chan iface.IWsSession
	LoginCh      chan iface.IWsSession
	UnRegisterCh chan iface.IWsSession
}

func GetSessionMgr() *SessionMgr {
	return sessionMgr
}

func NewSessionMgr() *SessionMgr {
	s := &SessionMgr{
		clients: make(map[iface.IWsSession]struct{}),
		online:  make(map[uint64]iface.IWsSession),

		RegisterCh:   make(chan iface.IWsSession, 1024),
		LoginCh:      make(chan iface.IWsSession, 1024),
		UnRegisterCh: make(chan iface.IWsSession, 1024),
	}

	return s
}

func (s *SessionMgr) AllLength() int {
	return len(s.clients)
}

func (s *SessionMgr) IsConn(ss iface.IWsSession) (ok bool) {
	s.clientsLock.RLock()
	defer s.clientsLock.RUnlock()

	_, ok = s.clients[ss]

	return
}

func (s *SessionMgr) GetAll() (allSS map[iface.IWsSession]struct{}) {
	allSS = make(map[iface.IWsSession]struct{})

	s.AllRange(func(ss iface.IWsSession) (result bool) {
		allSS[ss] = struct{}{}
		return true
	})

	return
}

func (s *SessionMgr) AllRange(fn func(ss iface.IWsSession) (result bool)) {
	s.clientsLock.RLock()
	defer s.clientsLock.RUnlock()

	for k := range s.clients {
		result := fn(k)
		if !result {
			return
		}
	}
}

func (s *SessionMgr) AllAdd(ss iface.IWsSession) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	s.clients[ss] = struct{}{}
}

func (s *SessionMgr) AllDel(ss iface.IWsSession) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	delete(s.clients, ss)
}

func (s *SessionMgr) OnlineLen() int {
	return len(s.online)
}

func (s *SessionMgr) OnlineAdd(userID uint64, ss iface.IWsSession) {
	s.onlineLock.Lock()
	defer s.onlineLock.Unlock()

	s.online[userID] = ss
}

func (s *SessionMgr) OnlineDel(userID uint64) {
	s.onlineLock.Lock()
	defer s.onlineLock.Unlock()

	delete(s.online, userID)
}

func (s *SessionMgr) OnlineGetOne(userID uint64) (ss iface.IWsSession) {
	s.onlineLock.RLock()
	defer s.onlineLock.RUnlock()

	if v, ok := s.online[userID]; ok {
		ss = v
	}

	return
}

func (s *SessionMgr) OnlineOnce(userID uint64, fn func(ss iface.IWsSession)) {
	cli := s.OnlineGetOne(userID)
	if cli == nil {
		return
	}

	fn(cli)
}

func (s *SessionMgr) OnlineRange(fn func(userID uint64, ss iface.IWsSession)) {
	s.onlineLock.RLock()
	defer s.onlineLock.RUnlock()

	for userID, ss := range s.online {
		fn(userID, ss)
	}

	return
}

func (s *SessionMgr) IsOnline(UID uint64) (ss iface.IWsSession, ok bool) {
	s.onlineLock.RLock()
	defer s.onlineLock.RUnlock()

	ss, ok = s.online[UID]
	return
}

func (s *SessionMgr) Login(ss iface.IWsSession) {
	if s.IsConn(ss) {
		s.OnlineAdd(ss.GetID(), ss)
	}
}
func (s *SessionMgr) Disconnect(ss iface.IWsSession) {
	if ss.GetID() != 0 {
		s.OnlineDel(ss.GetID())
	}

	s.AllDel(ss)
}

func (s *SessionMgr) Start() {
	for {
		select {
		case ss := <-s.RegisterCh:
			s.AllAdd(ss)
		case ss := <-s.LoginCh:
			s.Login(ss)
		case ss := <-s.UnRegisterCh:
			s.Disconnect(ss)
		}
	}
}

func (s *SessionMgr) ClearTimeout() {
	currentTime := time.Now()

	allClients := s.GetAll()
	for session := range allClients {
		var fn = func(args ...any) bool {
			return session.IsHeartbeatTimeout(currentTime)
		}
		if session.CheckSomething(fn) {
			session.Close()
		}
	}
}
