package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warning(msg string, args ...any)
	Error(msg string, args ...any)
}

type SLogger struct{}

func NewSLogger(level string) Logger {
	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: nil,
	}
	switch strings.ToLower(level) {
	case "debug":
		logConfig.Level = slog.LevelDebug
	case "info":
		logConfig.Level = slog.LevelInfo
	case "warning":
		logConfig.Level = slog.LevelWarn
	case "error":
		logConfig.Level = slog.LevelError
	}

	logHandler := slog.NewTextHandler(os.Stderr, logConfig)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	slog.Info("Logger started", "level", level)
	return &SLogger{}
}

func (*SLogger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func (*SLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (*SLogger) Warning(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (*SLogger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
