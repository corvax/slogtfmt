package slogtfmt

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler

	// AddSource causes the handler to compute the source code position
	// of the log statement and add a SourceKey attribute to the output.
	AddSource bool

	// TimeFormat is the format used for timestamps in the output.
	// If empty, the Handler will omit the timestamps.
	TimeFormat string

	// TimeInUTC specifies whether the time format should use UTC instead of the local time zone.
	TimeInUTC bool

	// TimeAttributeFormat specifies the time format used for the time attribute
	// in the log record. If empty, the default time format of time.RFC3339 is used.
	TimeAttributeFormat string

	// TimeAttributeInUTC specifies whether the time attribute in the log record
	// should use UTC instead of the local time zone.
	TimeAttributeInUTC bool
}

// Handler is a custom implementation of [slog.Handler] that provides advanced formatting capabilities
// for log records. It offers the following features:
//   - Customizable time value formatting for both log timestamps and time attributes
//   - Support for log record tagging using square brackets before the message
//   - Optional inclusion of source code information (file and line number)
type Handler struct {
	opts Options
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
}

type groupOrAttrs struct {
	attrs []slog.Attr
	group string
}

// tagKeyName is the key used to set the tag name attribute on a log record.
// The tag key value will be put in square brackets before the log message.
const tagKeyName = "__tag__"

// Tag returns an slog.Attr that can be used to set the tag for a log record.
// The tag value will be put in square brackets before the log message.
func Tag(name string) slog.Attr {
	return slog.Attr{Key: tagKeyName, Value: slog.StringValue(name)}
}

// NewHandler creates a new Handler with the provided io.Writer and Options.
// If no Options are provided, it will use the default Options with the time format set to time.RFC3339.
func NewHandler(out io.Writer, opts *Options) *Handler {
	h := &Handler{
		mu:  &sync.Mutex{},
		out: out,
	}
	if opts == nil {
		opts = &Options{
			TimeFormat:          time.RFC3339,
			TimeAttributeFormat: time.RFC3339,
		}
	}

	h.opts = *opts

	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}

	if h.opts.TimeAttributeFormat == "" {
		h.opts.TimeAttributeFormat = time.RFC3339
	}

	return h
}

// Enabled returns whether the given log level is enabled for this Handler.
// The Handler will only log records with a level greater than or equal to the configured level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle processes a log record and writes it to the configured io.Writer.
// It appends the time, level, tag (if set), source location (if configured),
// message, and attributes to the output. The output is formatted according to the
// configured Options.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	bufp := allocBuf()
	buf := *bufp
	defer func() {
		*bufp = buf
		freeBuf(bufp)
	}()

	// Append the time.
	if !r.Time.IsZero() && h.opts.TimeFormat != "" {
		if h.opts.TimeInUTC {
			buf = append(buf, r.Time.UTC().Format(h.opts.TimeFormat)...)
		} else {
			buf = append(buf, r.Time.Format(h.opts.TimeFormat)...)
		}
		buf = append(buf, "\t"...)
	}

	// Append the level.
	buf = append(buf, r.Level.String()...)

	goas := h.goas
	// Append the tag. Tag must be set by With().
	for _, goa := range goas {
		for _, a := range goa.attrs {
			if a.Key == tagKeyName {
				buf = append(buf, "\t["...)
				buf = append(buf, a.Value.String()...)
				buf = append(buf, "]"...)
				break
			}
		}
	}

	// Append the source.
	if h.opts.AddSource {
		frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()

		buf = append(buf, "\t"...)
		buf = append(buf, frame.File...)
		buf = append(buf, ":"...)
		buf = strconv.AppendInt(buf, int64(frame.Line), 10)
	}

	// Append the message.
	buf = append(buf, "\t"...)
	buf = append(buf, r.Message...)

	// Append the groups.
	if r.NumAttrs() == 0 {
		// If the record has no Attrs, remove groups at the end of the list
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}
	groupPrefix := ""
	for _, goa := range goas {
		if goa.group != "" {
			groupPrefix += goa.group + "."
		}
		for _, a := range goa.attrs {
			if a.Key != tagKeyName {
				buf = h.appendAttr(buf, a, groupPrefix)
			}
		}
	}

	// Append the attributes.
	r.Attrs(func(attr slog.Attr) bool {
		buf = h.appendAttr(buf, attr, groupPrefix)
		return true
	})

	buf = append(buf, "\n"...)

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

// WithGroup returns a new Handler that will log all records with the given group name.
// If the group name is empty, the original Handler is returned.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

// WithAttrs returns a new Handler that will log all records with the given attributes.
// If the list of attributes is empty, the original Handler is returned.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

// appendAttr appends the given attribute to the provided buffer, with the given prefix.
// It handles different attribute value types, including strings, times, and attribute groups.
// Attributes with empty values are ignored.
func (h *Handler) appendAttr(buf []byte, attr slog.Attr, prefix string) []byte {
	// Resolve the Attr's value before doing anything else.
	attr.Value = attr.Value.Resolve()

	// Ignore empty attrs.
	if attr.Equal(slog.Attr{}) {
		return buf
	}

	// Separate the attributes with a space.
	buf = append(buf, " "...)

	switch attr.Value.Kind() {
	case slog.KindString:
		buf = append(buf, prefix+attr.Key...)
		buf = append(buf, "="...)
		buf = strconv.AppendQuote(buf, attr.Value.String())
	case slog.KindTime:
		buf = append(buf, prefix+attr.Key...)
		buf = append(buf, "="...)
		if h.opts.TimeAttributeInUTC {
			buf = append(buf, attr.Value.Time().UTC().Format(h.opts.TimeAttributeFormat)...)
		} else {
			buf = append(buf, attr.Value.Time().Format(h.opts.TimeAttributeFormat)...)
		}
	case slog.KindGroup:
		attrs := attr.Value.Group()

		// Ignore empty groups.
		if len(attrs) == 0 {
			return buf
		}

		// If the Key is not empty, write it out.
		if attr.Key != "" {
			prefix = prefix + attr.Key + "."
		}

		for _, a := range attrs {
			buf = h.appendAttr(buf, a, prefix)
		}
	default:
		buf = append(buf, prefix+attr.Key...)
		buf = append(buf, "="...)
		buf = append(buf, attr.Value.String()...)
	}
	return buf
}

// withGroupOrAttrs creates a new Handler with the provided groupOrAttrs added to the list of goas.
// This allows the Handler to be configured with additional groups or attributes to be included
// in the formatted log output.
func (h *Handler) withGroupOrAttrs(goa groupOrAttrs) *Handler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}
