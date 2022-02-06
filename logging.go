package commons

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// +-------------------+
// | Stateless logging |
// +-------------------+

func Code(e *zerolog.Event, code string) {
	e.Str("code", code).Str("caller", Location3()).Msg(Location2())
}

func What(e *zerolog.Event, what string) {
	e.Str("what", what).Str("caller", Location3()).Msg(Location2())
}

func Msg(e *zerolog.Event) {
	e.Str("caller", Location3()).Msg(Location2())
}

// +--------------+
// | FilterWriter |
// +--------------+

type FilterWriter struct {
	Level  zerolog.Level
	Writer io.Writer
}

func (w *FilterWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)

	if err != nil {
		err = fmt.Errorf("FilterWriter.Write: %w", err)
	}

	return
}

func (w *FilterWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.Level {
		n, err = w.Writer.Write(p)

		if err != nil {
			err = fmt.Errorf("FilterWriter.WriteLevel: %w", err)
		}

		return
	}

	return len(p), nil
}

// +-----------------------+
// | Logging configuration |
// +-----------------------+

var ErrUnrecognizedLogLevel = errors.New("unrecognized log level")

// https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d
type LogConfig struct {
	Level string

	// Enable console logging.
	ConsoleLoggingEnabled bool

	// Directory to log to to when filelogging is enabled.
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory.
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled.
	MaxSize int

	// MaxBackups the max number of rolled files to keep.
	MaxBackups int

	// MaxAge the max age in days to keep a logfile.
	MaxAge int
}

func parseLogLevel(level string) (zerolog.Level, error) {
	switch level {
	case "trace":
		return zerolog.TraceLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "panic":
		return zerolog.PanicLevel, nil
	case "no":
		return zerolog.NoLevel, nil
	case "disabled":
		return zerolog.Disabled, nil
	}

	return zerolog.Disabled, fmt.Errorf("%w: %s", ErrUnrecognizedLogLevel, level)
}

func SetupLogger(c LogConfig) error {
	logLevel, err := parseLogLevel(c.Level)
	if err != nil {
		return err
	}

	// Set desired time format.
	zerolog.DurationFieldInteger = false
	// zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldUnit = time.Second
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	// Set global debugging level.
	zerolog.SetGlobalLevel(logLevel)

	// Setup outputs.
	var writers []io.Writer

	if c.ConsoleLoggingEnabled {
		var console zerolog.ConsoleWriter
		console.Out = os.Stderr

		writers = append(writers, &FilterWriter{
			Level:  logLevel, // Or zerolog.InfoLevel?
			Writer: console,
		})
	}

	if c.Filename != "" {
		w := &FilterWriter{
			Level: zerolog.DebugLevel,
			Writer: &lumberjack.Logger{
				Filename:   path.Join(c.Directory, c.Filename),
				MaxSize:    c.MaxSize,
				MaxAge:     c.MaxAge,
				MaxBackups: c.MaxBackups,
				LocalTime:  false,
				Compress:   false,
			},
		}

		writers = append(writers, w)
	}

	multi := zerolog.MultiLevelWriter(writers...)

	// Create logger with all settings taken into account.
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	Msg(log.Debug().
		Str("filename", c.Filename).
		Int("maxSize", c.MaxSize).
		Int("maxAge", c.MaxAge).
		Int("maxBackups", c.MaxBackups))

	return nil
}
