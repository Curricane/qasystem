package db

import (
	"fmt"
	"github.com/Curricane/logger"
	"github.com/jmoiron/sqlx"
	"qasystem/common"
)

//参数检查
func checkAnswer(answer *common.Answer) (err error) {
	if answer == nil ||
		answer.AnswerId <= 0 ||
		answer.AuthorId <= 0 ||
		len(answer.Content) == 0 ||
		len(answer.QuestionId) == 0 ||
		answer.CommentCount < 0 ||
		answer.VoteupCount < 0 {
		err = fmt.Errorf("nil or invalid answer, answer is:%v", answer)
		return
	}
	return
}

// 创建answer
func CreateAnswer(answer *common.Answer, qid int64) (err error) {

	err = checkAnswer(answer)
	if err != nil {
		return
	}

	if qid < 0 {
		err = fmt.Errorf("qid < 0, qid: %d", qid)
		return
	}

	// step1 存储answer基本信息
	sqlstr := "insert into answer(answer_id, content, author_id) values(?, ?, ?)"
	// 涉及到一个业务多个表操作，用事务提交
	tx, err := DB.Begin()
	if err != nil {
		logger.Error("DB.Begin failed, err is:%v", err)
		return
	}

	_, err = tx.Exec(sqlstr, answer.AnswerId, answer.Content, answer.AuthorId)
	if err != nil {
		tx.Rollback()
		logger.Error("create answer failed, question:%#v, err:%v", answer, err)
		return
	}

	// step2 问题-回答 关系表
	sqlstr = "insert into question_answer_rel(question_id, answer_id)values(?, ?)"
	_, err = tx.Exec(sqlstr, qid, answer.AnswerId)
	if err != nil {
		tx.Rollback()
		logger.Error("insert into question_answer_rel failed, err:%v", err)
		return
	}

	tx.Commit()
	return
}

// 根据question_id 获取一组answer_id， offset偏移 limit 数量
func GetAnswerIdList(qid int64, offset, limit int64)(answerIdList []int64, err error) {
	if qid <= 0 || offset < 0 || limit < 0 {
		err = fmt.Errorf("invalid param, qid: %d, offset: %d, limit: %d", qid, offset, limit)
		return
	}

	sqlstr := "select answer_id from question_answer_rel where question_id=? order by id desc limit ?, ?"
	err = DB.Select(&answerIdList, sqlstr, qid, offset, limit)

	if err != nil {
		logger.Error("get answer list failed, err:%v", err)
		return
	}

	return
}

// 根据answerIds，获取一组answer
func MGetAnswer(answerIds []int64) (answerList []*common.Answer, err error) {
	if len(answerIds) == 0 {
		logger.Warn("len(answerIds) == 0")
	}

	sqlstr := `select 
				answer_id, content, comment_count,
				voteup_count, author_id, status, can_comment,
				create_time, update_time
			   from
				answer where answer_id in(?)`
	var interfaceSlice []interface{}
	for _, c := range answerIds {
		interfaceSlice = append(interfaceSlice, c)
	}
	// 把一个？扩展为in (?, ?, ?)等
	insqlStr, params, err := sqlx.In(sqlstr, interfaceSlice) // interfaceSlice不能展开
	if err != nil {
		logger.Error("sqlx.in failed, sqlstr:%v, err:%v", sqlstr, err)
		return
	}
	logger.Debug("insqlStr:%v", insqlStr)
	logger.Debug("params:%v", params)

	err = DB.Select(&answerList, insqlStr, params...) // params需要展开
	if err != nil {
		logger.Error("MGetAnswer  failed, insqlStr:%v, category_ids:%v, err:%v",
			insqlStr, answerIds, err)
		return
	}

	return
}

// question_id对应的回答数
func GetAnswerCount(qid int64) (answerCount int64, err error) {
	if qid <= 0 {
		logger.Error("qid <= 0, qid: %d", qid)
		err = fmt.Errorf("qid <= 0, qid: %d", qid)
		return
	}

	sqlstr := "select count(answer_id) from question_answer_rel where question_id=?"
	err = DB.Get(&answerCount, sqlstr, qid)
	if err != nil {
		logger.Error("get GetAnswerCount failed, err:%v", err)
		return
	}
	return
}

func UpdateAnswerLikeCount(answerId int64) (err error) {
	sqlstr := `update answer set voteup_count=voteup_count+1 where answer_id=?`
	_, err = DB.Exec(sqlstr, answerId)
	if err != nil {
		logger.Error("UpdateAnswerLikeCount failed, err:%v", err)
		return
	}

	return
}


