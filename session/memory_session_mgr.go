package session

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

type MemorySessionMgr struct {
	sessionMap map[string]Session
	rwlock     sync.RWMutex
}

// Init 在内存中，不需要addr和option参数
func (s *MemorySessionMgr) Init(addr string, option ...string) (err error) {
	return
}

func (s *MemorySessionMgr) Get(sessionid string) (session Session, err error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	session, ok := s.sessionMap[sessionid]
	if !ok {
		err = ErrSessionNotExit
		return
	}
	return
}

func (s *MemorySessionMgr) CreateSession() (session Session, err error) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	id := uuid.NewV4()

	sessionId := id.String()
	session = NewMemorySession(sessionId)

	s.sessionMap[sessionId] = session

	return
}

func NewMemorySession() *SessionMgr {
	sr := &MemorySessionMgr{
		sessionMap: make(map[string]Session, 1024),
	}
	return sr
}
