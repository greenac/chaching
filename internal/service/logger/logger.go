package logger

import (
	"github.com/rs/zerolog"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type logType int

const (
	logTypeDebug logType = iota
	logTypeInfo
	logTypeWarn
	logTypeError
)

type ILogger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

var _ ILogger = (*ZeroLogWrapper)(nil)

type ZeroLogWrapper struct {
	logger   zerolog.Logger
	logLevel LogLevel
}

func (l *ZeroLogWrapper) Info(msg string) {
	l.log(logTypeInfo, msg)
}

func (l *ZeroLogWrapper) Warn(msg string) {
	l.log(logTypeWarn, msg)
}

func (l *ZeroLogWrapper) Error(msg string) {
	l.log(logTypeError, msg)
}

func (l *ZeroLogWrapper) Debug(msg string) {
	l.log(logTypeDebug, msg)
}

func (l *ZeroLogWrapper) log(lt logType, msg string) {
	if !l.shouldLog(lt) {
		return
	}

	switch lt {
	case logTypeInfo:
		l.logger.Info().Msg(msg)
	case logTypeWarn:
		l.logger.Warn().Msg(msg)
	case logTypeError:
		l.logger.Error().Msg(msg)
	case logTypeDebug:
		l.logger.Debug().Msg(msg)
	}
}

func (l *ZeroLogWrapper) shouldLog(t logType) bool {
	return int(t) >= int(l.logLevel)
}

func LogLevelForLogLevelName(ll string) LogLevel {
	switch strings.ToLower(ll) {
	case "debug":
		return LogLevelDebug
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}
