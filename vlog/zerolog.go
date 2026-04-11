package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type zeroLoggerImpl struct {
	logger zerolog.Logger
	ctx    context.Context
}

func NewZeroLogger(app string, level Level) Logger {
	l := &zeroLoggerImpl{
		zerolog.New(os.Stdout).Level(level.toZerologLevel()).With().Timestamp().Stack().Logger(),
		context.Background(),
	}
	if app != "" {
		l.logger.Hook(AppHook(app))
	}
	return l
}

type zerologEvent struct {
	*zerolog.Event
}

func (c *zerologEvent) Any(key string, i any) Event {
	return &zerologEvent{c.Event.Any(key, i)}
}

func (l *zeroLoggerImpl) Trace() Event {
	return &zerologEvent{l.logger.Trace().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Debug() Event {
	return &zerologEvent{l.logger.Debug().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Info() Event {
	return &zerologEvent{l.logger.Info().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Warn() Event {
	return &zerologEvent{l.logger.Warn().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Error() Event {
	return &zerologEvent{l.logger.Error().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Fatal() Event {
	return &zerologEvent{l.logger.Fatal().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) Panic() Event {
	return &zerologEvent{l.logger.Panic().Ctx(l.ctx)}
}

func (l *zeroLoggerImpl) WithCtx(ctx context.Context) Logger {
	return &zeroLoggerImpl{logger: l.logger, ctx: ctx}
}
