package util

/*
后端与前端的错误码与消息
*/

// 返回给前端的错误码，增加ErrCode需要相应增加GetMessage
const (
	ErrCodeSuccess 		= 0
	ErrCodeParameter 	= 1001
	ErrCodeUserExist 	= 1002
	ErrCodeServerBusy 	= 1003
	ErrCodeUserNotExit	= 1004
	ErrCodeUserPasswordWrong = 1005
)

// 返回给前端的错误码对应点额错误消息
func GetMessage(code int) (msg string) {
	switch code {
	case ErrCodeSuccess:
		msg = "Success"
	case ErrCodeParameter:
		msg = "Invalid parameter"
	case ErrCodeUserExist:
		msg = "用户已经存在"
	case ErrCodeServerBusy:
		msg = "服务器繁忙 "
	case ErrCodeUserNotExit:
		msg = "用户不存在"
	case ErrCodeUserPasswordWrong:
		msg = "用户或密码错误"
	default:
		msg = "Unknow error"
	}
	return
}
