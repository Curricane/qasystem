package session

import "fmt"

var (
	sessionMgr SessionMgr
)

// Init provider 1 memory 2 redis
func Init(provider string, addr string, options ...string) (err error) {
	switch provider {
	case "memory":
		sessionMgr = NewMemorySessionMgr()
	case "redis":
		sessionMgr = NewRedisSessionMgr()
	default:
		err = fmt.Errorf("not support")
		return
	}
	err = sessionMgr.Init(addr, options...)
	return
}

func CreateSession() (session Session, err error) {
	return sessionMgr.CreateSession()
}
func Get(sessionid string) (session Session, err error) {
	return sessionMgr.Get(sessionid)
}


