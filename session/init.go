package session

import "fmt"

var (
	sessionMgr SessionMgr
)

// Init provider 1 memory 2 redis
func Init(provider string, addr string, options ...string) (err error) {
	switch provider {
	case "memory":
		sessionMgr = NewMemorySession()
	case "redis":
		sessionMgr = NewRedisSession()
	default:
		err = fmt.Errorf("not support")
		return
	}
	err = sessionMgr.Init(addr, options)
	return
}
