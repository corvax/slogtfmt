package extras

import (
	"log/slog"
	"os"

	"github.com/corvax/slogtfmt"
	slogexp "github.com/smallnest/slog-exp"
)

// NewSplitLevelLogHandler creates a new slog.Handler with the provided level.
// It will log messages to stdout for INFO, WARN, and DEBUG levels and to stderr for ERROR level.
func NewSplitLevelLogHandler(level slog.Level, opts *slogtfmt.Options) slog.Handler {
	// Use common options for all handlers.
	// Level in the mapped handler is not important (it will be ignored by the LevelHandler).
	if opts == nil {
		opts = &slogtfmt.Options{
			TimeFormat: slogtfmt.RFC3339Milli,
		}
	}

	// Add handlers for all available log levels.
	allLevelHandlers := map[slog.Level]slog.Handler{
		slog.LevelInfo:  slogtfmt.NewHandler(os.Stdout, opts),
		slog.LevelWarn:  slogtfmt.NewHandler(os.Stdout, opts),
		slog.LevelDebug: slogtfmt.NewHandler(os.Stdout, opts),
		slog.LevelError: slogtfmt.NewHandler(os.Stderr, opts),
	}

	// Remove the handlers that are below the specified level.
	// LevelHandler will simply call Handle() method of the handler mapped to a level.
	// So, the levels below the specified level should be removed from the map.
	for l := range allLevelHandlers {
		if l < level {
			delete(allLevelHandlers, l)
		}
	}

	return slogexp.NewLevelHandler(allLevelHandlers)
}

// NewSplitLevelLogger creates a new slog.Logger instance with the provided level.
// It will log messages to stdout for INFO, WARN, and DEBUG levels and to stderr for ERROR level.
func NewSplitLevelLogger(level slog.Level, opts *slogtfmt.Options) *slog.Logger {
	return slog.New(NewSplitLevelLogHandler(level, opts))
}
