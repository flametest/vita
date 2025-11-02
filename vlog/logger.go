package log

import "github.com/rs/zerolog"

type Logger interface {
	Debug(args ...interface{}) *zerolog.Event
	Info(args ...interface{}) *zerolog.Event
	Warn(args ...interface{}) *zerolog.Event
	Error(args ...interface{}) *zerolog.Event
	Fatal(args ...interface{}) *zerolog.Event
	Panic(args ...interface{}) *zerolog.Event
}

type loggerImpl struct {
	zerolog.Logger
}

func (l *loggerImpl) Debug(args ...interface{}) *zerolog.Event {
	return l.Logger.Debug()
}

func (l *loggerImpl) Info(args ...interface{}) *zerolog.Event {
	return l.Logger.Info()
}

func (l *loggerImpl) Warn(args ...interface{}) *zerolog.Event {
	return l.Logger.Warn()
}

func (l *loggerImpl) Error(args ...interface{}) *zerolog.Event {
	return l.Logger.Error()
}

func (l *loggerImpl) Fatal(args ...interface{}) *zerolog.Event {
	return l.Logger.Fatal()
}

func (l *loggerImpl) Panic(args ...interface{}) *zerolog.Event {
	return l.Logger.Panic()
}
