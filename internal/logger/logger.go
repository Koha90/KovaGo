// Package logger configures structed application logging.
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Koha90/KovaGo/internal/config"
)

// New creates a configured slog logger.
//
// The returned cleanup function releases resources owned by the logger,
// such as an opened log file.
func New(cfg config.Log) (*slog.Logger, func() error, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, nil, err
	}

	options := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	}

	consoleHandler, err := newHandler(
		cfg.Format,
		os.Stdout,
		options,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"create console log handler: %w",
			err,
		)
	}

	if !cfg.File.Enabled {
		return slog.New(consoleHandler), noopCleanup, nil
	}

	if strings.TrimSpace(cfg.File.Path) == "" {
		return nil, nil, fmt.Errorf(
			"log file path is required when file logging is enabled",
		)
	}

	if err := os.MkdirAll(
		filepath.Dir(cfg.File.Path),
		0o755,
	); err != nil {
		return nil, nil, fmt.Errorf(
			"create log directory: %w",
			err,
		)
	}

	file, err := os.OpenFile(
		cfg.File.Path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o640,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"open log file %q: %w",
			cfg.File.Path,
			err,
		)
	}

	fileHandler, err := newHandler(
		cfg.File.Format,
		file,
		options,
	)
	if err != nil {
		_ = file.Close()

		return nil, nil, fmt.Errorf(
			"create file log handler: %w",
			err,
		)
	}

	handler := slog.NewMultiHandler(
		consoleHandler,
		fileHandler,
	)

	return slog.New(handler), file.Close, nil
}

func newHandler(
	format string,
	writer io.Writer,
	options *slog.HandlerOptions,
) (slog.Handler, error) {
	switch strings.ToLower(format) {
	case "text":
		return slog.NewTextHandler(
			writer,
			options,
		), nil

	case "json":
		return slog.NewJSONHandler(
			writer,
			options,
		), nil

	default:
		return nil, fmt.Errorf(
			"unsupported log format %q",
			format,
		)
	}
}

func parseLevel(value string) (slog.Level, error) {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug, nil

	case "info":
		return slog.LevelInfo, nil

	case "warn", "warning":
		return slog.LevelWarn, nil

	case "error":
		return slog.LevelError, nil

	default:
		return slog.LevelInfo, fmt.Errorf(
			"unsupported log level %q",
			value,
		)
	}
}

func noopCleanup() error {
	return nil
}

// Discard creates a logger that discard all records.
//
// It is primarily useful in tests.
func Discard() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}
