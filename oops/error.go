package oops

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

var (
	ErrFormInvalid = errors.New("form: invalid")

	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrNotEnough = errors.New("models: balance not enough")

	ErrNotAllowed = errors.New("models: not allowed")
)

func ErrorLog(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
}

type Error struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) Error {
	return Error{Code: code, Msg: msg}
}