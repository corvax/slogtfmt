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

// NewLogger creates a new LoggerFormatter that wraps the provided slog.Logger.
// The LoggerFormatter provides formatted logging methods that delegate to the underlying slog.Logger.
func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger}
}

// Infof logs a formatted info message.
func (l *Logger) Infof(format string, args ...any) {
	if l.Logger.Enabled(context.Background(), slog.LevelInfo) {
		l.Logger.Info(fmt.Sprintf(format, args...))
	}
}

// InfofContext logs a formatted info message with context.
func (l *Logger) InfofContext(ctx context.Context, format string, args ...any) {
	if l.Logger.Enabled(ctx, slog.LevelInfo) {
		l.Logger.InfoContext(ctx, fmt.Sprintf(format, args...))
	}
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(format string, args ...any) {
	if l.Logger.Enabled(context.Background(), slog.LevelDebug) {
		l.Logger.Debug(fmt.Sprintf(format, args...))
	}
}

// DebugfContext logs a formatted debug message with context.
func (l *Logger) DebugfContext(ctx context.Context, format string, args ...any) {
	if l.Logger.Enabled(ctx, slog.LevelDebug) {
		l.Logger.DebugContext(ctx, fmt.Sprintf(format, args...))
	}
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...any) {
	if l.Logger.Enabled(context.Background(), slog.LevelWarn) {
		l.Logger.Warn(fmt.Sprintf(format, args...))
	}
}

// WarnfContext logs a formatted warning message with context.
func (l *Logger) WarnfContext(ctx context.Context, format string, args ...any) {
	if l.Logger.Enabled(ctx, slog.LevelWarn) {
		l.Logger.WarnContext(ctx, fmt.Sprintf(format, args...))
	}
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...any) {
	if l.Logger.Enabled(context.Background(), slog.LevelError) {
		l.Logger.Error(fmt.Sprintf(format, args...))
	}
}

// ErrorfContext logs a formatted error message with context.
func (l *Logger) ErrorfContext(ctx context.Context, format string, args ...any) {
	if l.Logger.Enabled(ctx, slog.LevelError) {
		l.Logger.ErrorContext(ctx, fmt.Sprintf(format, args...))
	}
}
