package log

import "context"

type BaseLogger interface {
	Trace() Event
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
	Panic() Event
}

type Logger interface {
	BaseLogger
	WithCtx(ctx context.Context) Logger
}
