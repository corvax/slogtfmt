package loggerf

import (
	"context"
	"fmt"
	"log/slog"
)

// Logger wraps slog.Logger and adds formatted logging methods.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new Logger that wraps the provided slog.Logger.
// The Logger provides formatted logging methods that delegate to the underlying slog.Logger.
func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger}
}

// Logf logs a formatted message at the specified log level.
func (l *Logger) Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	// Check if the logger is enabled to avoid unnecessary calls of fmt.Sprintf().
	if l.Logger.Enabled(ctx, level) {
		l.Logger.Log(ctx, level, fmt.Sprintf(format, args...))
	}
}

// Infof logs a formatted info message.
func (l *Logger) Infof(format string, args ...any) {
	l.Logf(context.Background(), slog.LevelInfo, format, args...)
}

// InfofContext logs a formatted info message with context.
func (l *Logger) InfofContext(ctx context.Context, format string, args ...any) {
	l.Logf(ctx, slog.LevelInfo, format, args...)
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(format string, args ...any) {
	l.Logf(context.Background(), slog.LevelDebug, format, args...)
}

// DebugfContext logs a formatted debug message with context.
func (l *Logger) DebugfContext(ctx context.Context, format string, args ...any) {
	l.Logf(ctx, slog.LevelDebug, format, args...)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...any) {
	l.Logf(context.Background(), slog.LevelWarn, format, args...)
}

// WarnfContext logs a formatted warning message with context.
func (l *Logger) WarnfContext(ctx context.Context, format string, args ...any) {
	l.Logf(ctx, slog.LevelWarn, format, args...)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...any) {
	l.Logf(context.Background(), slog.LevelError, format, args...)
}

// ErrorfContext logs a formatted error message with context.
func (l *Logger) ErrorfContext(ctx context.Context, format string, args ...any) {
	l.Logf(ctx, slog.LevelError, format, args...)
}
