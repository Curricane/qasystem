package account

import (
	"github.com/gin-gonic/gin"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/id_gen"
	"qasystem/util"
)

func LoginHandle(ctx *gin.Context) {
	// step1 获取登录信息 UserInfo
	var userInfo common.UserInfo
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Password) == 0 || len(userInfo.Username) == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	err = db.Login(&userInfo)
	switch err {
	case db.ErrUserNotExits:
		util.ResponseError(ctx, util.ErrCodeUserNotExit)
		return
	case db.ErrUserPasswordWrong:
		util.ResponseError(ctx, util.ErrCodeUserPasswordWrong)
		return
	case nil:
		break
	default:
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	// 登录成功
	util.ResponseSuccess(ctx, nil)
}

func RegisterHandle(ctx *gin.Context) {

	// step1 获取注册信息 UserInfo
	var userInfo common.UserInfo
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Email) == 0 || len(userInfo.Password) == 0 ||
		len(userInfo.Username) == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	if userInfo.Sex != common.UserSexMan && userInfo.Sex != common.UserSexWomen {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	userInfo.UserId, err = id_gen.GetId()
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	err = db.Register(&userInfo)
	if err == db.ErrUserExits {
		util.ResponseError(ctx, util.ErrCodeUserExist)
		return
	}

	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}
	util.ResponseSuccess(ctx, nil)
}


