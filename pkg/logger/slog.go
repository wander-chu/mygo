package logger

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func Debug(ctx context.Context, message string, args ...any) {
	log(ctx, slog.LevelDebug, message, args...)
}

func Info(ctx context.Context, message string, args ...any) {
	log(ctx, slog.LevelInfo, message, args...)
}

func Warn(ctx context.Context, message string, args ...any) {
	log(ctx, slog.LevelWarn, message, args...)
}

func Error(ctx context.Context, message string, args ...any) {
	log(ctx, slog.LevelError, message, args...)
}

func log(ctx context.Context, level slog.Level, message string, args ...any) {
	FromContext(ctx).Log(ctx, level, message, args...)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
