package main

import (
	"github.com/gin-gonic/gin"
)

func initTemplate(router *gin.Engine) {

	// 路径映射，访问/ 实际访问/static/index.html
	router.StaticFile("/", "/static/index.html")
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
	initTemplate(router)

	router.Run(":8080")
}
