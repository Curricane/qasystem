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
	logger.Debug("create question succ, question:%#v", question)

	err = db.CreateQuestion(&question)
	if err != nil {
		logger.Error("failed to store question into db, err is:%v", err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
	}

	util.ResponseSuccess(ctx, nil)

}
