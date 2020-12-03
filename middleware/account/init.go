package account

import (
	"qasystem/session"
)

// 账号验证前，需要初始化会话模块
func InitSession(provider string, addr string, options ...string) (err error) {
	return session.Init(provider, addr, options...)
}
