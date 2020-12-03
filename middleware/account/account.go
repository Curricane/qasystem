package account

import (
	"github.com/gin-gonic/gin"
	"qasystem/session"
)

// 请求处理前的账号中间件函数
func processRequest(ctx *gin.Context) {

	var userSession session.Session
	var err error

	// 最后确保有一个会话存到gin框架中
	defer func() {
		if userSession == nil {
			userSession, err = session.CreateSession()
		}

		ctx.Set(QAsystemSessionName, userSession)
	}()

	// step1 获取cookie
	cookie, err := ctx.Request.Cookie(CookieSessionId)
	if err != nil {
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step2 从cookie中获取sessionId
	sessionId := cookie.Value
	if len(sessionId) == 0 {
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step3 根据sessionId获取session（服务端）
	userSession, err =  session.Get(sessionId)
	if err != nil {
		// 获取不到，则认为没有登录过，设置用户登录状态0
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step4 获取已登录的用户id
	tmpUserId, err := userSession.Get(QAsystemUserId)
	if err != nil {
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}
	userId, ok := tmpUserId.(int64)
	if !ok || userId == 0 {
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step5 在gin中设置当前会话用户登录状态（已登录）
	ctx.Set(QAsystemUserId, int64(userId))
	ctx.Set(QAsystemUserLoginStatus, int64(1))

	return
}
