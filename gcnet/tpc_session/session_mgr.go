package tpc_session

import (
	"github.com/v587-zyf/gc/iface"
	"sync"
)

var sessionMgr *SessionMgr

func init() {
	sessionMgr = NewSessionMgr()
}

type SessionMgr struct {
	onlineMap sync.Map // uid:iface.ITcpSession
}

func GetSessionMgr() *SessionMgr {
	return sessionMgr
}

func NewSessionMgr() *SessionMgr {
	s := &SessionMgr{}

	return s
}

func (s *SessionMgr) Length() int {
	length := 0
	s.onlineMap.Range(func(key, value any) bool {
		length++
		return true
	})

	return length
}

func (s *SessionMgr) GetOne(UID uint64) iface.ITcpSession {
	if v, ok := s.onlineMap.Load(UID); ok {
		return v.(iface.ITcpSession)
	}
	return nil
}

func (s *SessionMgr) IsOnline(UID uint64) bool {
	_, ok := s.onlineMap.Load(UID)
	return ok
}

func (s *SessionMgr) Add(ss iface.ITcpSession) {
	s.onlineMap.Store(ss.GetID(), ss)
}

func (s *SessionMgr) Disconnect(SID uint64) {
	s.onlineMap.Delete(SID)
}

func (s *SessionMgr) Once(UID uint64, fn func(SS iface.ITcpSession)) {
	cli := s.GetOne(UID)
	if cli == nil {
		fn(nil)
		return
	}

	fn(cli)
}

func (s *SessionMgr) Range(fn func(UID uint64, SS iface.ITcpSession)) {
	s.onlineMap.Range(func(key, value interface{}) bool {
		uid := key.(uint64)
		ss := value.(iface.ITcpSession)
		fn(uid, ss)
		return true
	})
}

func (s *SessionMgr) Close() {
	s.onlineMap.Range(func(key, value interface{}) bool {
		uid := key.(uint64)
		ss := value.(iface.ITcpSession)
		ss.Close()
		s.onlineMap.Delete(uid)
		return true
	})
}
