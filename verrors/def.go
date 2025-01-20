package verrors

import "errors"

var (
	ErrTypeAssertion = errors.New("type assertion failed")
)
