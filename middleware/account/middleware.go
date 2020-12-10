package account

import (
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"qasystem/util"
)

func AuthMiddleware(ctx *gin.Context) {
	ProcessRequest(ctx)
	isLogin := IsLogin(ctx)
	if !isLogin {
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		logger.Debug("user is not logined")
		// 中断当前请求
		ctx.Abort()
		return
	}
	ctx.Next()
}
