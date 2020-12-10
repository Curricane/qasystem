package util

/*
后端与前端的错误码与消息
*/

// 返回给前端的错误码，增加ErrCode需要相应增加GetMessage
const (
	ErrCodeSuccess 			= 0
	ErrCodeParameter 		= 1001
	ErrCodeUserExist 		= 1002
	ErrCodeServerBusy 		= 1003
	ErrCodeUserNotExit	    = 1004
	ErrCodeUserPasswordWrong = 1005
	ErrCodeCaptionHit        = 1006
	ErrCodeContentHit        = 1007
	ErrCodeNotLogin          = 1008
	ErrCodeRecordExist       = 1009
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
	case ErrCodeCaptionHit:
		msg = "标题中含有非法内容, 请修改后发表"
	case ErrCodeContentHit:
		msg = "内容中含有非法内容，请修改后发表"
	case ErrCodeNotLogin:
		msg = "用户未登录"
	case ErrCodeRecordExist:
		msg = "记录已经存在"
	default:
		msg = "Unknow error"
	}
	return
}
