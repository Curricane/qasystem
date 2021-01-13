package favorite

import (
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/id_gen"
	"qasystem/middleware/account"
	"qasystem/util"
	"strconv"
	"strings"
)

const (
	MinCommentCotentSize = 10
)

/* 创建收藏夹
type FavoriteDir struct {
	DirId int64  // 无， 后端生成
	DirName string  // 前端提供
	Count   int32  // 无，首次创建为0
	UserId  int64  // 无，后端从session中获取
}
 */
func AddDirHandle(ctx *gin.Context) {
	var fvDir common.FavoriteDir
	err := ctx.BindJSON(&fvDir)
	if err != nil {
		logger.Error("failed to BindJSON common.FavoriteDir")
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}
	logger.Debug("bind json succ, favoriteDir:%#v", fvDir)

	dir_id, err := id_gen.GetId()
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("id_gen.GetId failed, favoriteDir:%#v, err:%v", fvDir, err)
		return
	}
	fvDir.DirId = int64(dir_id)

	userId, err := account.GetUserId(ctx)
	if err != nil || userId == 0 {
		logger.Error("failed GetUserId")
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	fvDir.UserId = userId
	err = db.CreateFavoriteDir(&fvDir)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("CreateFavorieDir failed, comment:%#v, err:%v", fvDir, err)
		return
	}

	util.ResponseSuccess(ctx, nil)
}

/*添加到收藏夹
type Favorite struct {
	AnswerId int64 // 前端提供
	UserId   int64 // 无，后端从session中获取
	DirId    int64 // 前端提供
}
 */
func AddFavoriteHandle(ctx *gin.Context) {
	var fv common.Favorite
	err := ctx.BindJSON(&fv)
	if err != nil {
		logger.Error("failed to BindJSON common.Favorite")
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}
	logger.Debug("bind json succ, favorite:%#v", fv)

	if fv.DirId == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("invalid favorite:%v", fv)
		return
	}

	userId, err := account.GetUserId(ctx)
	if err != nil || userId == 0 {
		logger.Error("failed GetUserId")
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	fv.UserId = userId

	err = db.CreateFavorite(&fv)
	switch err {
	case nil:
		break
	case db.ErrRecordExists:
		util.ResponseError(ctx, util.ErrCodeRecordExist)
		logger.Warn("fv had existed in db, favorite:%#v, err:%v", fv, err)
		return
	default:
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("CreateFavorie failed, favorite:%#v, err:%v", fv, err)
		return
	}

	util.ResponseSuccess(ctx, nil)
}

// 根据userId 获取收藏夹列表
func DirListHandle(ctx *gin.Context) {
	userId, err := account.GetUserId(ctx)
	if err != nil || userId == 0 {
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}

	fvDirList, err := db.GetFavoriteDirList(userId)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("GetFavoriteDirList failed, user_id:%v, err:%v", userId, err)
		return
	}

	util.ResponseSuccess(ctx, fvDirList)
}

func FavoriteListHandle(ctx *gin.Context) {
	dirIdStr, ok := ctx.GetQuery("dir_id")
	dirIdStr = strings.TrimSpace(dirIdStr)
	if ok == false || len(dirIdStr) == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid dir id, val:%v", dirIdStr)
		return
	}
	dirId, err := strconv.ParseInt(dirIdStr, 10, 64)
	if err != nil || dirId == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid dir id, val:%v", dirIdStr)
		return
	}
	logger.Debug("get query dir_id succ, val:%v", dirIdStr)

	//解析offset
	var offset int64
	offsetStr, ok := ctx.GetQuery("offset")
	offsetStr = strings.TrimSpace(offsetStr)
	if ok == false || len(offsetStr) == 0 {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	logger.Debug("get query offset succ, val:%v", offsetStr)

	//解析limit
	var limit int64
	limitStr, ok := ctx.GetQuery("limit")
	limitStr = strings.TrimSpace(limitStr)
	if ok == false || len(limitStr) == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}
	logger.Debug("get query limit succ, val:%v", limitStr)

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}
	logger.Debug("get query limit succ, val:%v", limitStr)

	// 获取userId
	userId, err := account.GetUserId(ctx)
	if err != nil || userId == 0 {
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}

	fvList, err := db.GetFavoriteList(userId, dirId, offset, limit)
	if err != nil {
		logger.Error("GetFavoriteList failed, dir_id:%v, user_id:%v, err:%v", userId, dirId, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	var answerIdList []int64
	for _, v := range fvList {
		answerIdList = append(answerIdList, v.AnswerId)
	}

	answerList, err := db.MGetAnswer(answerIdList)
	if err != nil {
		logger.Error("db.MGetAnswer failed, answer_ids:%v err:%v",
			answerIdList, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	var userIdList []int64
	for _, v := range answerList {
		userIdList = append(userIdList, v.AuthorId)
	}

	userInfoList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		logger.Error("db.GetUserInfoList failed, user_ids:%v err:%v",
			userIdList, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	apiAnswerList := &common.ApiAnswerList{}
	for _, v := range answerList {
		apiAnswer := &common.ApiAnswer{}
		apiAnswer.Answer = *v

		for _, user := range userInfoList {
			if int64(user.UserId) == v.AuthorId {
				apiAnswer.AuthorName = user.Username
				break
			}
		}

		apiAnswerList.AnswerList = append(apiAnswerList.AnswerList, apiAnswer)
	}

	util.ResponseSuccess(ctx, apiAnswerList)
}
