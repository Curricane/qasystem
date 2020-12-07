package util

/*
后端与前端的错误码与消息
*/

const (
	ErrCodeSuccess 		= 0
	ErrCodeParameter 	= 1001
	ErrCodeUserExist 	= 1002
	ErrCodeServerBusy 	= 1003
)

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
	default:
		msg = "Unknow error"
	}
	return
}
