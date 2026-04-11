package log

import "context"

type LogType string

const (
	ZerologType LogType = "zerolog"
)

var logger Logger

func InitLogger(logType LogType, service string, level Level) {
	switch logType {
	case ZerologType:
		logger = NewZeroLogger(service, level)
	default:
		logger = NewZeroLogger(service, level)
	}
}

func WithCtx(ctx context.Context) BaseLogger {
	return logger.WithCtx(ctx)
}

func Debug() Event {
	return logger.Debug()
}

func Info() Event {
	return logger.Info()
}

func Warn() Event {
	return logger.Warn()
}

func Error() Event {
	return logger.Error()
}

func Panic() Event {
	return logger.Panic()
}

func Fatal() Event {
	return logger.Fatal()
}
