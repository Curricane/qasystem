package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
给前端的返回值
{
	"code":0, // 0表示成功，其他表示失败
	"message": "success",  //返回值的描述
	"data": { // 返回给前端的数据

	}
}
*/
type ResponseData struct {
	Code 	int 					`json:"code"`
	Message string 					`json:"message"`
	Data 	interface{} 	`json:"data"`
}

// 根据错误码code返回错误信息
func ResponseError(ctx *gin.Context, code int) {
	responseData := &ResponseData{
		Code: code,
		Message: GetMessage(code),
		Data: nil,
	}

	ctx.JSON(http.StatusOK, responseData)
}

// 返回成功信息，并提交数据给前端
func ResponseSuccess(ctx *gin.Context, data interface{}) {
	responseData := &ResponseData{
		Code: ErrCodeSuccess,
		Message: GetMessage(ErrCodeSuccess),
		Data: data,
	}
	ctx.JSON(http.StatusOK, responseData)
}
