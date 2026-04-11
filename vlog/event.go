package log

type Event interface {
	Msg(string)
	Msgf(string, ...any)
	Any(key string, i any) Event
}
