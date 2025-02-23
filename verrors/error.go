package verrors

import (
	"fmt"
	"github.com/flametest/vita/vhttp"
)

type Error struct {
	err      error
	errCode  int
	errMsg   string
	httpCode vhttp.StatusCode
	service  string
}

func (e *Error) Err() error {
	return e.err
}

func (e *Error) ErrCode() int {
	return e.errCode
}

func (e *Error) ErrMsg() string {
	return e.errMsg
}

func (e *Error) HttpCode() vhttp.StatusCode {
	return e.httpCode
}

func (e *Error) Service() string {
	return e.service
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s|%d|%s]: %v]", e.service, e.errCode, e.errMsg, e.err.Error())
}
