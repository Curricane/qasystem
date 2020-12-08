package db

import "errors"

var (
	ErrUserExits = errors.New("user has existed")
	ErrUserNotExits = errors.New("user not exist")
	ErrUserPasswordWrong = errors.New("username or password not right")
)
