package vmiddleware

import (
	"fmt"

	"github.com/flametest/vita/verrors"
	"github.com/flametest/vita/vhttp"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func NewErrorResponse(e *verrors.Error) *vhttp.VResponse {
	return &vhttp.VResponse{
		Service: e.Service(),
		Code:    e.ErrCode(),
		Message: e.ErrMsg(),
		Error:   e.Error(),
		Stack:   fmt.Sprintf("%+v", errors.WithStack(e)),
	}
}

func ErrorHandleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			var newErr *verrors.Error
			if err != nil {
				if ve, ok := err.(*verrors.Error); ok {
					newErr = ve
				} else if ee, ok := err.(*echo.HTTPError); ok {
					newErr = verrors.NewFromEchoHTTPError(ee, ee.Code)
				} else {
					newErr = verrors.InternalServerError(err.Error())
				}
			}
			if newErr != nil {
				res := NewErrorResponse(newErr)
				return c.JSON(newErr.HttpCode().Int(), res)
			}
			return err
		}
	}
}
