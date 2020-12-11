package question

import (
	"github.com/gin-gonic/gin"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/filter"
	"qasystem/id_gen"
	"qasystem/middleware/account"
	"qasystem/util"
	"github.com/Curricane/logger"
	"strconv"
)

func QuestionSubmitHandle(ctx *gin.Context) {

	// #step1 构建问题信息结构
	var question common.Question
	err := ctx.BindJSON(&question)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}
	logger.Debug("bind json success, question.Caption:%#v", question.Caption)

	// #step2 过滤标题和文本的敏感词
	result, hit := filter.Replace(question.Caption, "***")
	if hit {
		logger.Error("caption is hit filter, result:%v", result)
		util.ResponseError(ctx, util.ErrCodeCaptionHit)
		return
	}
	result, hit = filter.Replace(question.Content, "***")
	if hit {
		logger.Error("Content is hit filter, result:%v", result)
		util.ResponseError(ctx, util.ErrCodeContentHit)
		return
	}
	logger.Debug("filter succ, result:%#v\", result")

	// #step3 生成唯一id
	qid, err := id_gen.GetId()
	if err != nil {
		logger.Error("generate question id failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}
	question.QuestionId = int64(qid)

	// #step4 绑定作者信息
	userId, err := account.GetUserId(ctx)
	if err != nil || userId <= 0{
		logger.Error("user is not login, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	question.AuthorId = userId
	logger.Debug("create question succ, question:%#v", question)

	err = db.CreateQuestion(&question)
	if err != nil {
		logger.Error("failed to store question into db, err is:%v", err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
	}

	util.ResponseSuccess(ctx, nil)

}

// 获取问题详情
func QuestionDetailHandle(ctx *gin.Context) {
	questionIdStr, ok := ctx.GetQuery("question_id")
	if !ok {
		logger.Error("invalid question_id, not found question_id")
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		logger.Error("invalid question_id, strconv.ParseInt failed, err:%v, str:%v",
			err, questionIdStr)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	question, err := db.GetQuestion(questionId)
	if err != nil {
		logger.Error("get question failed, err:%v, str:%v", err, questionIdStr)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	categoryMap, err := db.MGetCategory([]int64{question.CategoryId})
	if err != nil {
		logger.Error("get category failed, err:%v, question:%v", err, question)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	category, ok := categoryMap[question.CategoryId]
	if !ok {
		logger.Error("get category failed, err:%v, question:%v", err, question)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	userInfoList, err := db.GetUserInfoList([]int64{question.AuthorId})
	if err != nil || len(userInfoList) == 0 {
		logger.Error("get user info list failed,user_ids:%#v, err:%v",
			question.AuthorId, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	apiQuestionDetail := &common.ApiQuestionDetail{}
	apiQuestionDetail.Question = *question
	apiQuestionDetail.AuthorName = userInfoList[0].Username
	apiQuestionDetail.CategoryName = category.CategoryName

	util.ResponseSuccess(ctx, apiQuestionDetail)
}
