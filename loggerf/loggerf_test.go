package loggerf

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/corvax/slogtfmt"
	"github.com/stretchr/testify/assert"
)

func TestLoggerf(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		message  string
		args     []any
		expected string
	}{
		{"Debug message", slog.LevelDebug, "Debug message: %s", []any{"Hello, World!"}, "DEBUG\tDebug message: Hello, World!\n"},
		{"Info message", slog.LevelInfo, "User %s logged in from %s", []any{"user", "localhost"}, "INFO\tUser user logged in from localhost\n"},
		{"Warn message", slog.LevelWarn, "Warning: disk usage is at %d%%", []any{98}, "WARN\tWarning: disk usage is at 98%\n"},
		{"Error message", slog.LevelError, "Error occured: %v", []any{fmt.Errorf("test error")}, "ERROR\tError occured: test error\n"},
	}

	var buf bytes.Buffer
	handler := slogtfmt.NewHandler(&buf, &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	})
	slogger := slog.New(handler)
	logger := NewLogger(slogger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			logger.Logf(context.Background(), tt.level, tt.message, tt.args...)
			// t.Log(buf.String())
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestLoggerf_With(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		message  string
		args     []any
		expected string
	}{
		{"Debug message", slog.LevelDebug, "Debug message: %s", []any{"Hello, World!"}, "DEBUG\t[test]\tDebug message: Hello, World!\n"},
		{"Info message", slog.LevelInfo, "User %s logged in from %s", []any{"user", "localhost"}, "INFO\t[test]\tUser user logged in from localhost\n"},
		{"Warn message", slog.LevelWarn, "Warning: disk usage is at %d%%", []any{98}, "WARN\t[test]\tWarning: disk usage is at 98%\n"},
		{"Error message", slog.LevelError, "Error occured: %v", []any{fmt.Errorf("test error")}, "ERROR\t[test]\tError occured: test error\n"},
	}

	var buf bytes.Buffer
	handler := slogtfmt.NewHandler(&buf, &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	})
	slogger := slog.New(handler)
	logger := NewLogger(slogger.With(slogtfmt.Tag("test")))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			logger.Logf(context.Background(), tt.level, tt.message, tt.args...)
			// t.Log(buf.String())
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestLoggerf_Args(t *testing.T) {
	var buf bytes.Buffer
	handler := slogtfmt.NewHandler(&buf, &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	})
	slogger := slog.New(handler)
	logger := NewLogger(slogger)

	expected := "INFO\tUser is logged in username=\"user\" host=\"localhost\"\n"
	username := "user"
	host := "localhost"
	logger.Info("User is logged in", "username", username, "host", host)
	// t.Log(buf.String())
	assert.Equal(t, expected, buf.String())
}
