package logger

import (
	"github.com/rs/zerolog"
)

type ILogger interface {
	Level(lvl zerolog.Level) zerolog.Logger
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
}

type Logger struct {
	Logger ILogger
}
