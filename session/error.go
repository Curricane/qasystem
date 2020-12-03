package session

import (
	"errors"
)

var (
	ErrSessionNotExit = errors.New("session not exists")
	ErrKeyNotExistInSession = errors.New("key not exists in session")
)