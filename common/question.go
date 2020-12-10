package common

import "time"

type Question struct {
	QuestionId int64 `json:"question_id_number", db:"question_id"`
	Caption       string    `json:"caption" db:"caption"`
	Content       string    `json:"content" db:"content"`
	AuthorId      int64     `json:"author_id_number" db:"author_id"`
	CategoryId    int64     `json:"category_id" db:"category_id"`
	Status        int32     `json:"status" db:"status"`
	CreateTime    time.Time `json:"-" db:"create_time"`
	CreateTimeStr string    `json:"create_time"` // 转为str用于给前端 "2006/1/2 15:04:05"
	QuestionIdStr string    `json:"question_id"` // 转为str用于给前端
	AuthorIdStr   string    `json:"author_id"` // 转为str用于给前端
}

// 用于给前端展示的结构
type ApiQuestion struct {
	Question
	AuthorName string `json:"author_name"`
}

type ApiQuestionDetail struct {
	Question
	AuthorName	string	`json:"author_name"`
	CategoryName	string	`json:"category_name"`
}
