package vmiddleware

import (
	"github.com/flametest/vita/vhttp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TraceWithRequestId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := middleware.DefaultRequestIDConfig
			req := c.Request()
			res := c.Response()

			rid := req.Header.Get(vhttp.HeaderXRequestID)
			if rid == "" {
				rid = cfg.Generator()
				req.Header.Set(vhttp.HeaderXRequestID, rid)
				res.Header().Set(vhttp.HeaderXRequestID, rid)
			}
			return next(c)
		}
	}
}
