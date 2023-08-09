package main

import (
	"fmt"
	"log"
	"time"

	"log/slog"

	"github.com/getsentry/sentry-go"
	slogsentry "github.com/samber/slog-sentry"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://xxxxxxx@yyyyyyy.ingest.sentry.io/zzzzzzz",
		EnableTracing: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer sentry.Flush(2 * time.Second)

	logger := slog.New(slogsentry.Option{Level: slog.LevelDebug}.NewSentryHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("a message")
}
