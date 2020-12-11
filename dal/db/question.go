package db

import (
	"github.com/Curricane/logger"
	"qasystem/common"
)

func CreateQuestion(q *common.Question) (err error) {
	sqlstr := "insert into question (question_id, caption, content, author_id, category_id) values(?,?,?,?,?)"
	_, err = DB.Exec(sqlstr, q.QuestionId, q.Caption, q.Content, q.AuthorId, q.CategoryId)
	if err != nil {
		logger.Error("create question failed, question:%#v, err:%v", q, err)
		return
	}
	return
}

func GetQuestion(qid int64) (q *common.Question, err error) {
	q = &common.Question{}
	sqlstr := `select question_id, caption, content, author_id, category_id, create_time from question where question_id = ?`
	err = DB.Get(q, sqlstr, qid)
	if err != nil {
		logger.Error("get question  failed, err:%v", err)
		return
	}
	return
}

func GetQuestionList(cid int64) (questionList []*common.Question, err error) {
	sqlstr :=`select question_id, caption, content, author_id, category_id, create_time from question where category_id = ?`
	err = DB.Select(&questionList, sqlstr, cid)
	if err != nil {
		logger.Error("get question list failed, err:%v", err)
		return
	}
	return
}