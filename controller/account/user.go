package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginViewHandle(ctx *gin.Context) {
	// 使用时需要在html模板中，定义好名字 {{ define "views/login.html" }} {{ end }}
	ctx.HTML(http.StatusOK, "views/login.html", nil)
}

func RegisterViewHandle(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/register.html", nil)
}


