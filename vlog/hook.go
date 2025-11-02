package vlog

import "github.com/rs/zerolog"

type app string

func (h app) Run(e *zerolog.Event, level zerolog.Level, message string) {
	e.Str("app", string(h))
}

func AppHook(name string) zerolog.Hook {
	return app(name)
}
