package vmiddleware

import (
	"regexp"
	"time"

	log "github.com/flametest/vita/vlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func APILogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		cfg := middleware.DefaultLoggerConfig
		re, _ := regexp.Compile("health-check")
		cfg.Skipper = func(context echo.Context) bool {
			if re.Match([]byte(context.Request().URL.Path)) {
				return true
			}
			return false
		}
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)
			log.Info("duration: ", duration)
			return err
		}
	}
}
