package slogsentry

import (
	"log/slog"

	"github.com/getsentry/sentry-go"
)

var LogLevels = map[slog.Level]sentry.Level{
	slog.LevelDebug: sentry.LevelDebug,
	slog.LevelInfo:  sentry.LevelInfo,
	slog.LevelWarn:  sentry.LevelWarning,
	slog.LevelError: sentry.LevelError,
}
