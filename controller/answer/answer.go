package answer

import (
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"html"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/id_gen"
	"qasystem/middleware/account"
	"qasystem/util"
	"strconv"
)

// 根据question_id 获取一组回答
func AnswerListHandle(ctx *gin.Context) {
	// step1 从gin获取question_id offset limit
	qid, err := util.GetQueryInt64(ctx, "question_id")
	if err != nil {
		logger.Error("get question id failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	offset, err := util.GetQueryInt64(ctx, "offset")
	if err != nil {
		logger.Error("get offset failed, set 0 as default, err:%v", err)
		offset = 0
	}

	limit, err := util.GetQueryInt64(ctx, "limit")
	if err != nil {
		logger.Error("get limit failed, set 100 as default, err:%v", err)
		limit = 0
	}
	logger.Debug("get answer list parameter succ, qid:%v, offset:%v, limit:%v",
		qid, offset, limit)

	// step2 从数据库中获取一组answerids
	answerIdList, err := db.GetAnswerIdList(qid, offset, limit)
	if err != nil {
		logger.Error("db.GetAnswerIdList failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}
	if len(answerIdList) == 0 {
		logger.Warn("len(answerIdList) == 0")
		util.ResponseSuccess(ctx, "")
		return
	}

	// step3 根据answerids 获取一组answers
	answerList, err := db.MGetAnswer(answerIdList)
	if err != nil {
		logger.Error("db.MGetAnswer failed, answer_ids:%v err:%v",
			answerIdList, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	// step4 获取回答者的userinfo，用于组ApiAnswer
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

	// step5 组装ApiAnswer
	apiAnswerList := &common.ApiAnswerList{}
	for _, v := range answerList {
		apiAnswer := &common.ApiAnswer{}
		apiAnswer.Answer = *v

		for _, user := range userInfoList{
			if int64(user.UserId) == v.AuthorId {
				apiAnswer.AuthorName = user.Username
				break
			}
		}
		apiAnswerList.AnswerList = append(apiAnswerList.AnswerList, apiAnswer)
	}

	// 获取问题数
	count, err := db.GetAnswerCount(qid)
	if err != nil {
		logger.Error("db.GetAnswerCount failed, question_id:%v err:%v",
			qid, err)
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		return
	}

	apiAnswerList.TotalCount = int32(count)

	// step6 向前端返回结果
	util.ResponseSuccess(ctx, apiAnswerList)
}

// 提交回答
func AnswerCommitHandle(ctx *gin.Context) {
	// step1 从gin上获取answer关键字段
	var answer common.Answer
	err := ctx.BindJSON(&answer)
	if err != nil {
		logger.Error("bind json failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	// step2 获取question_id 用于关联answer
	qid, err := strconv.ParseInt(answer.QuestionId, 10, 64)
	if err != nil {
		logger.Error("invalid question_id")
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	// step3 获取userid，关联用户信息
	uid, err := account.GetUserId(ctx)
	if err != nil || uid == 0 {
		logger.Error("get user id failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	answer.AuthorId = uid
	// 1. 针对content做一个转义，防治xss漏洞
	answer.Content = html.EscapeString(answer.Content)
	// 2. 生成唯一评论id
	cid, err := id_gen.GetId()
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("id_gen.GetId failed, comment:%#v, err:%v", answer, err)
		return
	}
	answer.AnswerId = int64(cid)

	// step4 存入数据库中 answer库 + question_answer_rel库
	err = db.CreateAnswer(&answer, qid)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("CreatePostComment failed, comment:%#v, err:%v", answer, err)
		return
	}

	// step5 反馈给前端
	util.ResponseSuccess(ctx, nil)
}
