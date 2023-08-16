package utils

import (
	"context"
	"github.com/greenac/chaching/internal/service/logger"
)

func AddLoggerToCtx(ctx context.Context, baseLogger logger.ILogger, props map[string]string) context.Context {
	return context.WithValue(ctx, "logger", baseLogger.SubLogger(props))
}

func LoggerFromCtx(ctx context.Context) logger.ILogger {
	l, ok := ctx.Value("logger").(logger.ILogger)
	if !ok {
		l = logger.NewLogger(logger.LogLevelInfo, true)
	}

	return l
}
