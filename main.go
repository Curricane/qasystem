package main

import (
	"github.com/gin-gonic/gin"
	"qasystem/controller/account"
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

	err := initDb()
	if err != nil {
		panic(err)
	}

	// 初始化 id_gen，值后续需要更为实际的机器序列
	err = id_gen.Init(1)
	if err != nil {
		panic(err)
	}

	initTemplate(router)
	router.POST("/api/user/register", account.RegisterHandle)
	router.POST("/api/user/login", account.LoginHandle)

	router.Run(":9090")
}
