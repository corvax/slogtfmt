# slogtfmt

`slogtfmt` is a handler for slog that allows you to customize timestamp formats for both log timestamps and time attributes. This package also provides a helper function to add tags to log entries.

## Features

- Customizable time value formatting for both log timestamps and time attributes
- Support for log record tagging, tags are shown in square brackets before the message
- Optional inclusion of source code information (file and line number)

## Installation

To install `slogtfmt`, use `go get`:

```bash
go get github.com/corvax/slogtfmt
```
## Usage

### Creating a New Handler

To create a new slog handler, use the `NewHandler` function. This function requires an `io.Writer` where the logs will be written and an `Options` struct to configure the handler.

```go
package main

import (
	"time"
	"os"
	"github.com/corvax/slogtfmt"
)

func main() {
	opts := &slogtfmt.Options{
		Level:               slog.LevelInfo,
		AddSource:           false,
		TimeFormat:          "2006-01-02T15:04:05.999Z07:00",
		TimeInUTC:           true,
		TimeAttributeFormat: time.RFC3339,
		TimeAttributeInUTC:  true,
	}

	logger := slog.New(slogtfmt.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	slog.Info("Started", slog.Time("time", time.Now()))

	// To create a logger with an added tag, use With(slogtfmt.Tag("tag_name")
	serviceLogger := logger.With(slogtfmt.Tag("service"))
	serviceLogger.Info("Started")
}
```

#### Sample output

The output of the sample code above would be as following:

```text
2024-07-01T04:23:30.557Z    INFO    Started time=2024-07-01T03:41:05Z
2024-07-01T04:23:30.557Z    INFO    [service]     Started
```

### Without a timestamp

In order to omit the log timestamp, set `TimeFormat` to an empty string.

```go
	opts := &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	}

	slog.SetDefault(slog.New(slogtfmt.NewHandler(os.Stdout, opts)))

	slog.Info("Started", slog.Time("time", time.Now())
	slog.Debug("User connected", slog.String("user", "username"))
	slog.Warn("Access denied", slog.String("role", "readOnly"))
 	slog.Error("Service is unavailable")
 	slog.Info("Finished")
```

#### Output:

Note that the time attribute format is not affected and uses the default formatting `slogtfmt.RFC3339Milli`.
```
INFO	Started time=2024-07-01T15:02:50.720+10:00
DEBUG	User connected user="username"
WARN	Access denied role="readOnly"
ERROR	Service is unavailable
INFO	Finished
```

### Using `With` Option functions

You can also use the `slogtfmt.NewHandlerWithOptions()` constructor with Option functions.

To achieve the same log formatting as shown above, you can use the following snippet:

```go
slog.SetDefault(slog.New(slogtfmt.NewHandlerWithOptions(
	os.Stdout,
	slogtfmt.WithLevel(slog.LevelDebug),
	slogtfmt.WithTimeFormat(""),
)))
```

`With` functions are available for all `Options`.

### Default options

The constructor `slogtfmt.NewHandlerWithOptions()` creates the handler with the default `Options` and then updates them using the provided `With` option functions.

If `nil` as `Options` is passed to `slogtfmt.NewHandler()`, the handler will be created with the default options.

##### Default options

```go
defaultOptions := &Options{
	Level:               slog.LevelInfo,
	AddSource:           false,
	TimeFormat:          slogtfmt.RFC3339Milli,
	TimeInUTC:           false,
	TimeAttributeFormat: slogtfmt.RFC3339Milli,
	TimeAttributeInUTC:  false,
}
```

## Time formats

In addition to the standard time format, there are some additional time formats available in the package that can be used for formatting timestamps.

```go
const (
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"
)

```

## Options

The `Options` struct allows you to customize the behavior of the handler. Below is a detailed explanation of each field in the `Options` struct:

* **`Level`**: Specifies the minimum level to log. Logs with a lower level are discarded. If `nil`, the handler uses `slog.LevelInfo`.
* **`AddSource`**: If set to `true`, the handler computes the source code position of the log statement and adds the file name and the position to the output.
* **`TimeFormat`**: The format used for log timestamps in the output. If empty, the handler will omit the timestamps.
* **`TimeInUTC`**: Specifies whether the time format should use UTC instead of the local time zone.
* **`TimeAttributeFormat`**: Specifies the time format used for the time attribute in the log record. If empty, the default time format of `time.RFC3339` is used.
* **`TimeAttributeInUTC`**: Specifies whether the time attribute in the log record should use UTC instead of the local time zone.

## `loggerf.Logger`

The `loggerf.Logger` is a wrapper around the standard `slog.Logger` structure that provides additional functions to format log messages. These functions allow you to log messages with formatted strings, similar to `fmt.Printf`.

### Formatting functions

The `loggerf.Logger` provides the following functions to format the log message:

- `Debugf(format string, args ...any)`: Log a debug message with formatting.
- `Infof(format string, args ...any)`: Log an info message with formatting.
- `Warnf(format string, args ...any)`: Log a warning message with formatting.
- `Errorf(format string, args ...any)`: Log an error message with formatting.
- `Logf(ctx context.Context, level slog.Level, format string, args ...any)`: Log a formatted message at the specified log level.

### Usage

To use the `loggerf.Logger`, add the `loggerf` import to your code:

```bash
go get github.com/corvax/slogtfmt/loggerf
```

Create a new instance of `loggerf.Logger` based on any `slog.Logger`:

```go
import (
	"errors"
	"log/slog"
	"os"

	"github.com/corvax/slogtfmt"
	"github.com/corvax/slogtfmt/loggerf"
)

func main() {
	handler := slogtfmt.NewHandler(os.Stdout, &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	})
	slogger := slog.New(handler)
	logger := loggerf.NewLogger(slogger)

	// Now you can use the functions to format the log messages.
	// For example:
	username := "user"
	ipAddress := "localhost"
	diskUsage := 98
	err := errors.New("test error")

	logger.Debugf("Debug message: %s", "Hello, World!")
	logger.Infof("User %s logged in from %s", username, ipAddress)
	logger.Warnf("Warning: disk usage is at %d%%", diskUsage)
	logger.Errorf("Error occurred: %v", err)
}
```

Output:
```
DEBUG   Debug message: Hello, World!
INFO    User user logged in from localhost
WARN    Warning: disk usage is at 98%
ERROR   Error occured: test error
```

### With context

You can also create a new logger with additional context using the `With` method:

```go
import (
	"errors"
	"log/slog"
	"os"

	"github.com/corvax/slogtfmt"
	"github.com/corvax/slogtfmt/loggerf"
)

func main() {
	handler := slogtfmt.NewHandler(os.Stdout, &slogtfmt.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "",
	})
	slogger := slog.New(handler)
	logger := loggerf.NewLogger(slogger.With(slogtfmt.Tag("service")))

	// Now you can use the functions to format the log message.
	// For example:
	username := "user"
	ipAddress := "localhost"
	diskUsage := 98
	err := errors.New("test error")

	logger.Debugf("Debug message: %s", "Hello, World!")
	logger.Infof("User %s logged in from %s", username, ipAddress)
	logger.Warnf("Warning: disk usage is at %d%%", diskUsage)
	logger.Errorf("Error occurred: %v", err)
}
```

Output:
```
DEBUG   [service]	Debug message: Hello, World!
INFO    [service]	User user logged in from localhost
WARN    [service]	Warning: disk usage is at 98%
ERROR   [service]	Error occured: test error
```

`WithGroup` can also be used to create a new logger with additional context and a group name.

### Notes

Please be aware that the formatting log message functions use `args` to format the message. `args` values are not used to build `slog.Attr`.

If you want to add additional attributes to your log message, you can still use the `slog.Logger` methods, for example:

```go
logger.Info("User is logged in", "username", username, "host", ipAddress)
```

Output:
```
INFO    User is logged in username="user" host="localhost"
```

`loggerf.Logger` embeds the `slog.Logger` structure, so you can use any of the `slog.Logger` methods.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
