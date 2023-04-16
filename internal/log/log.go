package log

import (
	"github.com/hauntedness/giom/internal/runtime"
	"golang.org/x/exp/slog"
)

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// wrap Info with source file and line number
func Infos(args ...any) {
	slog.Info(runtime.SourceWithSkip(2), args...)
}

func Errors(err error, args ...any) {
	args = append([]any{"error", err}, args...)
	slog.Error(runtime.SourceWithSkip(2), args...)
}
