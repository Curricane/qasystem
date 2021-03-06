package session

// Session 会话接口
type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Save() error
	IsModify() bool
	Id() string
}
