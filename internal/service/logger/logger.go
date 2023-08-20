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
	InfoFmt(msg string, args ...any)
	Warn(msg string)
	WarnFmt(msg string, args ...any)
	Error(msg string)
	ErrorFmt(msg string, args ...any)
	Debug(msg string)
	DebugFmt(msg string, args ...any)
	SubLogger(props map[string]string) ILogger
}

func NewZeroLogWrapper(l zerolog.Logger, level LogLevel) ILogger {
	return &ZeroLogWrapper{logger: l, logLevel: level}
}

var _ ILogger = (*ZeroLogWrapper)(nil)

type ZeroLogWrapper struct {
	logger   zerolog.Logger
	logLevel LogLevel
}

func (l *ZeroLogWrapper) Info(msg string) {
	l.log(logTypeInfo, msg)
}

func (l *ZeroLogWrapper) InfoFmt(msg string, args ...any) {
	l.logFmt(logTypeInfo, msg, args...)
}

func (l *ZeroLogWrapper) Warn(msg string) {
	l.log(logTypeWarn, msg)
}

func (l *ZeroLogWrapper) WarnFmt(msg string, args ...any) {
	l.logFmt(logTypeWarn, msg, args...)
}

func (l *ZeroLogWrapper) Error(msg string) {
	l.log(logTypeError, msg)
}

func (l *ZeroLogWrapper) ErrorFmt(msg string, args ...any) {
	l.logFmt(logTypeError, msg, args...)
}

func (l *ZeroLogWrapper) Debug(msg string) {
	l.log(logTypeDebug, msg)
}

func (l *ZeroLogWrapper) DebugFmt(msg string, args ...any) {
	l.logFmt(logTypeDebug, msg, args...)
}

func (l *ZeroLogWrapper) SubLogger(props map[string]string) ILogger {
	log := l.logger
	for k, v := range props {
		log = log.With().Str(k, v).Logger()
	}

	return NewZeroLogWrapper(log, l.logLevel)
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

func (l *ZeroLogWrapper) logFmt(lt logType, msg string, args ...any) {
	if !l.shouldLog(lt) {
		return
	}

	switch lt {
	case logTypeInfo:
		l.logger.Info().Msgf(msg, args...)
	case logTypeWarn:
		l.logger.Warn().Msgf(msg, args...)
	case logTypeError:
		l.logger.Error().Msgf(msg, args...)
	case logTypeDebug:
		l.logger.Debug().Msgf(msg, args...)
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
