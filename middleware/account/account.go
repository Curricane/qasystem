package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qasystem/session"
)

// 请求处理前的账号中间件函数
func ProcessRequest(ctx *gin.Context) {

	var userSession session.Session
	var err error

	// 最后确保有一个会话存到gin框架中
	defer func() {
		if userSession == nil {
			fmt.Println("userSession == nil, CreateSession")
			userSession, err = session.CreateSession()
		}

		ctx.Set(QAsystemSessionName, userSession)
	}()

	// step1 获取cookie
	cookie, err := ctx.Request.Cookie(CookieSessionId)
	if err != nil {
		fmt.Println("no CookieSessionId in Cookie")
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step2 从cookie中获取sessionId
	sessionId := cookie.Value
	if len(sessionId) == 0 {
		fmt.Println("no CookieSessionId value in Cookie")
		ctx.Set(QAsystemUserId, int64(0))
		ctx.Set(QAsystemUserLoginStatus, int64(0))
		return
	}

	// step3 根据sessionId获取session（服务端）
	userSession, err =  session.Get(sessionId)
	if err != nil {
		// 获取不到，则认为没有登录过，设置用户登录状态0
		fmt.Println("no SessionId in Cookie")
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


func GetUserId(ctx *gin.Context) (userId int64, err error) {
	tmpUserId, ok := ctx.Get(QAsystemUserId)
	if !ok {
		err =  errors.New("userId not exists")
		return
	}

	userId, ok = tmpUserId.(int64)
	if !ok {
		err = errors.New("userId should be int64")
		return
	}

	return
}

// 设置用户id到session中
func SetUserId(userId int64, ctx *gin.Context) (err error) {
	// step1 获取当前Session
	var userSession session.Session
	tmpSession, ok := ctx.Get(QAsystemSessionName)
	if !ok {
		// 无session，不处理
		fmt.Println("cannt get tmpSession in ctx")
		return
	}
	userSession, ok = tmpSession.(session.Session)
	if !ok {
		// 错误的值，不处理
		fmt.Println("error type session")
		return
	}

	if userSession == nil {
		fmt.Println("userSession == nil")
		return
	}

	// 登录 首次设置cookie
	err = userSession.Set(QAsystemUserId, userId)
	return
}

func IsLogin(ctx *gin.Context)(login bool) {
	tmpLoginStatus, ok := ctx.Get(QAsystemUserLoginStatus)
	if !ok {
		return
	}

	loginStatus, ok := tmpLoginStatus.(int64)
	if !ok {
		return
	}

	if loginStatus == 0 {
		return
	}
	login = true
	return
}

// 响应处理前的账号中间件函数，看是否需要更新Session和Cookie
func ProcessResponse(ctx *gin.Context) {

	// step1 获取当前Session
	var userSession session.Session
	tmpSession, ok := ctx.Get(QAsystemSessionName)
	if !ok {
		// 无session，不处理
		fmt.Println("cannt get tmpSession in ctx")
		return
	}
	userSession, ok = tmpSession.(session.Session)
	if !ok {
		// 错误的值，不处理
		fmt.Println("error type session")
		return
	}

	// step2 session有修改，则在服务器里更新session
	if userSession.IsModify() == false {
		fmt.Println("userSession.IsModify() == false")
		return
	}
	err := userSession.Save()
	if err != nil {
		// 系统更新不了session，不继续处理
		fmt.Println("cannt save session")
		return
	}

	// session有更改，相应的更新cookie，通知客户端更新
	sessionId := userSession.Id()
	cookie := &http.Cookie{
		Name : CookieSessionId,
		Value: sessionId,
		MaxAge: CookieMaxAge,
		HttpOnly: true,
		Path: "/",
	}
	http.SetCookie(ctx.Writer, cookie)
	fmt.Println("set cookie ok")
	return
}
