package vlog

import (
	"os"

	"github.com/rs/zerolog"
)

var logger Logger

func InitLogger(app string, level zerolog.Level) {
	l := &loggerImpl{
		zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger(),
	}
	if app != "" {
		l.Hook(AppHook(app))
	}
	logger = l
}
