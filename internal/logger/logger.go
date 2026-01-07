package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	Log = slog.New(handler)
}

// Info logs info messages
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Error logs error messages
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

// Warn logs warning messages
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}
