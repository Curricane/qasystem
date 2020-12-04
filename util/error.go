package util

const (
	ErrCodeSuccess = 0
	ErrCodeParameter = 1001
)

func GetMessage(code int) (msg string) {
	switch code {
	case ErrCodeSuccess:
		msg = "Success"
	case ErrCodeParameter:
		msg = "Invalid parameter"
	default:
		msg = "Unknow error"
	}
	return
}
