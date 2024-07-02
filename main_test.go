package slogtfmt

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, &Options{
		TimeFormat: "",
	})
	logger := slog.New(handler)

	logger.Info("test message", "key1", "value1", "key2", 42)

	expected := "INFO\ttest message key1=\"value1\" key2=42\n"
	assert.Equal(t, expected, buf.String())

	buf.Reset()

	logger.Warn("warning message", "error", "something went wrong")

	expected = "WARN\twarning message error=\"something went wrong\"\n"
	assert.Equal(t, expected, buf.String())
}

func TestHandlerWithTags(t *testing.T) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, &Options{
		TimeFormat: "",
	})
	logger := slog.New(handler).With(Tag("my-tag"))

	logger.Info("test message", "key1", "value1", "key2", 42)

	expected := "INFO\t[my-tag]\ttest message key1=\"value1\" key2=42\n"
	assert.Equal(t, expected, buf.String())

	buf.Reset()

	logger.Warn("warning message", "error", "something went wrong")

	expected = "WARN\t[my-tag]\twarning message error=\"something went wrong\"\n"
	assert.Equal(t, expected, buf.String())
}

func BenchmarkHandler(b *testing.B) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, &Options{
		TimeFormat: "",
	})
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		logger.Info("benchmark message",
			"key1", "value1",
			"key2", true,
			"key3", 42,
			"key4", 3.14,
			"key5", time.Minute+time.Second,
			"key6", time.Now(), // 3 allocations when using slog.Time() attribute
		)
	}
}
