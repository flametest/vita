package vmiddleware

import (
	"fmt"

	"github.com/flametest/vita/verrors"
	"github.com/flametest/vita/vhttp"
	log "github.com/flametest/vita/vlog"
	"github.com/labstack/echo/v4"
)

func NewErrorResponse(e *verrors.Error, withStack bool) *vhttp.VResponse {
	resp := &vhttp.VResponse{
		Service: e.Service(),
		Code:    e.ErrCode(),
		Message: e.ErrMsg(),
		Error:   e.Error(),
	}
	if withStack {
		resp.Stack = fmt.Sprintf("%+v", verrors.WithStack(e))
	}
	return resp
}

func ErrorHandleMiddleware(withStack bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			log.Error().Err(err)
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
				res := NewErrorResponse(newErr, withStack)
				return c.JSON(newErr.HttpCode().Int(), res)
			}
			return err
		}
	}
}
