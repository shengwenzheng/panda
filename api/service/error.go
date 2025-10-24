package service

import "errors"

var (
	errSystem = errors.New("system error")
)

var ErrorCode = map[error]int{
	errSystem: 1000,
}