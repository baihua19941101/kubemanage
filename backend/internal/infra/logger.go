package infra

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"kubeManage/backend/internal/config"
)

func SetupLogger(cfg config.Config) (func() error, error) {
	level := parseLogLevel(cfg.LogLevel)
	format := strings.ToLower(strings.TrimSpace(cfg.LogFormat))
	output := strings.TrimSpace(cfg.LogOutput)
	if output == "" {
		output = "stdout"
	}

	writer, closer, err := buildLogWriter(output)
	if err != nil {
		return nil, err
	}

	handlerOpts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}
	slog.SetDefault(slog.New(handler))
	return closer, nil
}

func parseLogLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func buildLogWriter(output string) (io.Writer, func() error, error) {
	switch strings.ToLower(output) {
	case "stdout":
		return os.Stdout, func() error { return nil }, nil
	case "stderr":
		return os.Stderr, func() error { return nil }, nil
	default:
		dir := filepath.Dir(output)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, nil, fmt.Errorf("create log dir failed: %w", err)
			}
		}
		f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("open log file failed: %w", err)
		}
		return f, f.Close, nil
	}
}
