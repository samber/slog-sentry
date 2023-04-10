package slogsentry

import (
	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slog"
)

var levelMap = map[slog.Level]sentry.Level{
	slog.LevelDebug: sentry.LevelDebug,
	slog.LevelInfo:  sentry.LevelInfo,
	slog.LevelWarn:  sentry.LevelWarning,
	slog.LevelError: sentry.LevelError,
}
