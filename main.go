package main

import (
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"qasystem/controller/account"
	"qasystem/controller/answer"
	"qasystem/controller/category"
	"qasystem/controller/question"
	"qasystem/filter"
	mdlacc "qasystem/middleware/account"
	"qasystem/dal/db"
	"qasystem/id_gen"
)

func initDb() (err error) {
	dns := "root:cmc123456@tcp(localhost:3306)/qasystem?parseTime=true"
	err = db.Init(dns)
	if err != nil {
		return
	}
	return
}

func initSession() (err error) {
	err = mdlacc.InitSession("memory", "")
	return
}

func initFilter() (err error) {
	err = filter.Init("./data/filter.dat.txt")
	if err != nil {
		logger.Error("init filter failed, err:%v", err)
		return
	}

	logger.Debug("init filter succ")
	return
}

func initTemplate(router *gin.Engine) {

	// 路径映射，访问/ 实际访问/static/index.html
	router.StaticFile("/", "./static/index.html")
	//
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	// css静态资源路径映射
	router.Static("/css/", "./static/css/")
	// fonts静态资源路径映射
	router.Static("/fonts/", "./static/fonts/")
	// img静态资源路径映射
	router.Static("/img/", "./static/img/")
	// js静态资源路径映射
	router.Static("/js/", "./static/js/")
}

func main() {
	router := gin.Default()

	// 初始化日志模块
	config := make(map[string]string)
	config["log_level"] = "debug"

	logger.InitLogger("console", config)

	// 初始化DB模块
	err := initDb()
	if err != nil {
		panic(err)
	}

	// 初始化会话模块
	err = initSession()
	if err != nil {
		panic(err)
	}

	// 初始化 id_gen，值后续需要更为实际的机器序列
	err = id_gen.Init(1)
	if err != nil {
		panic(err)
	}

	err = initFilter()
	if err != nil {
		panic(err)
	}

	// 初始化前端资源
	initTemplate(router)

	// 设置路由
	router.POST("/api/user/register", account.RegisterHandle)
	router.POST("/api/user/login", account.LoginHandle)
	router.GET("/api/category/list", category.GetCategoryListHandle)
	router.POST("/api/ask/submit", mdlacc.AuthMiddleware, question.QuestionSubmitHandle)
	router.GET("/api/question/list", category.GetQuestionListHandle)
	router.GET("/api/question/detail", question.QuestionDetailHandle)
	router.GET("/api/answer/list", answer.AnswerListHandle)
	router.POST("/api/answer/commit", mdlacc.AuthMiddleware, answer.AnswerCommitHandle)
	router.Run(":9090")
}
