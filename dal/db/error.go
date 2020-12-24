package db

import "errors"

var (
	ErrUserExits = errors.New("username has existed")
	ErrUserNotExits = errors.New("username not exist")
	ErrUserPasswordWrong = errors.New("username or password not right")
	ErrRecordExists      = errors.New("record exist")
)
