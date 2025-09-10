package verrors

import (
	"fmt"
	"github.com/flametest/vita/vhttp"
	"github.com/pkg/errors"
	"net/http"
)

var service string

const DefaultErrorCode = 1000

type Error struct {
	err      error
	errCode  int
	errMsg   string
	httpCode vhttp.StatusCode
	service  string
}

func Initialize(serviceName string) {
	service = serviceName
	for _, predefinedError := range predefinedErrors {
		predefinedError.service = serviceName
	}
}

// New creates a new Error instance
func New(errCode int, errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  errCode,
		errMsg:   errMsg,
		httpCode: http.StatusInternalServerError,
		service:  service,
	}
}

func InternalServerError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusInternalServerError,
		service:  service,
	}
}

func BadRequestError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusBadRequest,
		service:  service,
	}
}

func NotFoundError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusNotFound,
		service:  service,
	}
}

func ForbiddenError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusForbidden,
		service:  service,
	}
}

func UnauthorizedError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusUnauthorized,
		service:  service,
	}
}

func NotImplementedError(errMsg string) *Error {
	return &Error{
		err:      errors.New(errMsg),
		errCode:  DefaultErrorCode,
		errMsg:   errMsg,
		httpCode: http.StatusNotImplemented,
		service:  service,
	}
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

func (e *Error) InternalServerError() *Error {
	e.httpCode = http.StatusInternalServerError
	return e
}
func (e *Error) BadRequest() *Error {
	e.httpCode = http.StatusBadRequest
	return e
}
func (e *Error) Unauthorized() *Error {
	e.httpCode = http.StatusUnauthorized
	return e
}

func (e *Error) Forbidden() *Error {
	e.httpCode = http.StatusForbidden
	return e
}
func (e *Error) NotFound() *Error {
	e.httpCode = http.StatusNotFound
	return e
}

func (e *Error) NotImplemented() *Error {
	e.httpCode = http.StatusNotImplemented
	return e
}

func Wrap(err error, a ...interface{}) error {
	return errors.Wrap(err, fmt.Sprint(a...))
}

func Wrapf(err error, format string, a ...interface{}) error {
	return errors.Wrapf(err, format, a...)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
