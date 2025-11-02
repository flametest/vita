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

func Debug(args ...interface{}) *zerolog.Event {
	return logger.Debug(args)
}

func Info(args ...interface{}) *zerolog.Event {
	return logger.Info(args)
}

func Warn(args ...interface{}) *zerolog.Event {
	return logger.Warn(args)
}

func Error(args ...interface{}) *zerolog.Event {
	return logger.Error(args)
}

func Panic(args ...interface{}) *zerolog.Event {
	return logger.Panic(args)
}

func Fatal(args ...interface{}) *zerolog.Event {
	return logger.Fatal(args)
}
