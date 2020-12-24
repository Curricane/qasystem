package category

import (
	"fmt"
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/util"
	"strconv"
)

func GetCategoryListHandle(ctx *gin.Context) {
	categoryList, err := db.GetCategoryList()
	if err != nil {
		logger.Error("db.GetCategoryList failed，err:%v", err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}
	util.ResponseSuccess(ctx, categoryList)
}


// 通过category_id 获取问题列表，问题列表中的问题是*common.ApiQuestion
func GetQuestionListHandle(ctx *gin.Context) {
	categoryIdStr, ok := ctx.GetQuery("category_id")
	if !ok {
		logger.Error("invalid category_id, not found category_id")
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		logger.Error("invalid category_id, strconv.ParseInt failed, err:%v, str:%v",
			err, categoryIdStr)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	questionList, err := db.GetQuestionList(categoryId)
	if err != nil {
		logger.Error("get question list failed,category_id:%v, err:%v",
			categoryId, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	if len(questionList) == 0 {
		logger.Warn("get question list succ, empty list,category_id:%v",
			categoryId)
		util.ResponseSuccess(ctx, questionList)
		return
	}

	// 获取questionList，所有的AuthorId
	var userIdList []int64
	userIdMap := make(map[int64]bool, 16)
	for _, question := range questionList {
		_, ok := userIdMap[question.AuthorId]
		if ok {
			continue
		}

		userIdMap[question.AuthorId] = true
		userIdList = append(userIdList, question.AuthorId)
	}

	// 通过userids获取用户信息
	userInfoList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		logger.Error("get user info list failed,user_ids:%#v, err:%v",
			userIdList, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}
	if len(userInfoList) == 0 {
		logger.Warn("len(userInfoList) == 0 ")
	}

	// 问题列表+用户信息 组成ApiQuestion
	var apiQuestionList []*common.ApiQuestion
	for _, q := range questionList {
		var apiQ = &common.ApiQuestion{}
		apiQ.Question = *q
		apiQ.QuestionIdStr = fmt.Sprintf("%d", apiQ.QuestionId)
		apiQ.AuthorIdStr = fmt.Sprintf("%d", apiQ.AuthorId)
		apiQ.CreateTimeStr = apiQ.CreateTime.Format("2006/1/2 15:04:05")
		logger.Info("apiQ.QuestionId is:%d, apiQ.QuestionIdStr:%v", apiQ.QuestionId, apiQ.QuestionIdStr)

		// 问题对应的 apiQ.AuthorName
		for _, userInfo := range userInfoList {
			logger.Debug("q.AuthorIdStr:%v, userInfo.UserId:%v, int64(userInfo.UserId):%v", q.AuthorIdStr, userInfo.UserId, int64(userInfo.UserId))
			if q.AuthorId == int64(userInfo.UserId) {
				apiQ.AuthorName = userInfo.Username
				break
 			}
		}

		apiQuestionList = append(apiQuestionList, apiQ)
	}

	// 组成ApiQuestion成功并返回给前端
	util.ResponseSuccess(ctx, apiQuestionList)
}