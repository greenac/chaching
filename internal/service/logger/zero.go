package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func NewLogger(logLevel LogLevel, shouldLogJson bool) ILogger {
	var logger zerolog.Logger
	if shouldLogJson {
		logger = zerolog.New(os.Stdout).With().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return &ZeroLogWrapper{
		logger:   logger,
		logLevel: logLevel,
	}
}
