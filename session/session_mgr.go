package session

// SessionMgr 会话管理 1:n Session
type SessionMgr interface {
	Init(addr string, options ...string) (err error)
	CreateSession() (session Session, err error)
	Get(sessionid string) (session Session, err error)
}
