package main

import (
	"github.com/gin-gonic/gin"
	accountctr "qasystem/controller/account"
)

func main() {
	router := gin.Default()

	router.Static("/static/", "./static")
	router.LoadHTMLGlob("views/*")

	router.GET("/user/login", accountctr.LoginViewHandle)
	router.GET("/user/register", accountctr.RegisterViewHandle)
	router.Run(":8080")
}
