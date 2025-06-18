
# slog: Sentry handler

[![tag](https://img.shields.io/github/tag/samber/slog-sentry.svg)](https://github.com/samber/slog-sentry/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/slog-sentry?status.svg)](https://pkg.go.dev/github.com/samber/slog-sentry)
![Build Status](https://github.com/samber/slog-sentry/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/slog-sentry)](https://goreportcard.com/report/github.com/samber/slog-sentry)
[![Coverage](https://img.shields.io/codecov/c/github/samber/slog-sentry)](https://codecov.io/gh/samber/slog-sentry)
[![Contributors](https://img.shields.io/github/contributors/samber/slog-sentry)](https://github.com/samber/slog-sentry/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/slog-sentry)](./LICENSE)

A [Sentry](https://sentry.io) Handler for [slog](https://pkg.go.dev/log/slog) Go library.

<div align="center">
  <hr>
  <sup><b>Sponsored by:</b></sup>
  <br>
  <a href="https://quickwit.io?utm_campaign=github_sponsorship&utm_medium=referral&utm_content=samber-slog-sentry&utm_source=github">
    <div>
      <img src="https://github.com/samber/oops/assets/2951285/49aaaa2b-b8c6-4f21-909f-c12577bb6a2e" width="240" alt="Quickwit">
    </div>
    <div>
      Cloud-native search engine for observability - An OSS alternative to Splunk, Elasticsearch, Loki, and Tempo.
    </div>
  </a>
  <hr>
</div>

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): `slog.Handler` chaining, fanout, routing, failover, load balancing...
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-sampling](https://github.com/samber/slog-sampling): `slog` sampling policy
- [slog-mock](https://github.com/samber/slog-mock): `slog.Handler` for test purposes

**HTTP middlewares:**

- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-echo](https://github.com/samber/slog-echo): Echo middleware for `slog` logger
- [slog-fiber](https://github.com/samber/slog-fiber): Fiber middleware for `slog` logger
- [slog-chi](https://github.com/samber/slog-chi): Chi middleware for `slog` logger
- [slog-http](https://github.com/samber/slog-http): `net/http` middleware for `slog` logger

**Loggers:**

- [slog-zap](https://github.com/samber/slog-zap): A `slog` handler for `Zap`
- [slog-zerolog](https://github.com/samber/slog-zerolog): A `slog` handler for `Zerolog`
- [slog-logrus](https://github.com/samber/slog-logrus): A `slog` handler for `Logrus`

**Log sinks:**

- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-betterstack](https://github.com/samber/slog-betterstack): A `slog` handler for `Betterstack`
- [slog-rollbar](https://github.com/samber/slog-rollbar): A `slog` handler for `Rollbar`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-sentry](https://github.com/samber/slog-sentry): A `slog` handler for `Sentry`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`
- [slog-quickwit](https://github.com/samber/slog-quickwit): A `slog` handler for `Quickwit`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-telegram](https://github.com/samber/slog-telegram): A `slog` handler for `Telegram`
- [slog-mattermost](https://github.com/samber/slog-mattermost): A `slog` handler for `Mattermost`
- [slog-microsoft-teams](https://github.com/samber/slog-microsoft-teams): A `slog` handler for `Microsoft Teams`
- [slog-webhook](https://github.com/samber/slog-webhook): A `slog` handler for `Webhook`
- [slog-kafka](https://github.com/samber/slog-kafka): A `slog` handler for `Kafka`
- [slog-nats](https://github.com/samber/slog-nats): A `slog` handler for `NATS`
- [slog-parquet](https://github.com/samber/slog-parquet): A `slog` handler for `Parquet` + `Object Storage`
- [slog-channel](https://github.com/samber/slog-channel): A `slog` handler for Go channels

## 🚀 Install

```sh
go get github.com/samber/slog-sentry/v2
```

**Compatibility**: go >= 1.21

No breaking changes will be made to exported APIs before v3.0.0.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/slog-sentry/v2](https://pkg.go.dev/github.com/samber/slog-sentry/v2)

### Handler options

```go
type Option struct {
    // Level sets the minimum log level to capture and send to Sentry.
    // Logs at this level and above will be processed. The default level is debug.
    Level slog.Leveler
    // Hub specifies the Sentry Hub to use for capturing events.
    // If not provided, the current Hub is used by default.
    Hub *sentry.Hub

    // Converter is an optional function that customizes how log records
    // are converted into Sentry events. By default, the DefaultConverter is used.
    Converter Converter
    // AttrFromContext is an optional slice of functions that extract attributes
    // from the context. These functions can add additional metadata to the log entry.
    AttrFromContext []func(ctx context.Context) []slog.Attr

    // AddSource is an optional flag that, when set to true, includes the source
    // information (such as file and line number) in the Sentry event.
    // This can be useful for debugging purposes.
    AddSource bool
    // ReplaceAttr is an optional function that allows for the modification or
    // replacement of attributes in the log record. This can be used to filter
    // or transform attributes before they are sent to Sentry.
    ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

    // BeforeSend is an optional function that allows for the modification of
    // the Sentry event before it is sent to the server. This can be used to add
    // additional context or modify the event payload.
    BeforeSend func(event *sentry.Event) *sentry.Event
}
```

Other global parameters:

```go
slogsentry.SourceKey = "source"
slogsentry.ContextKey = "extra"
slogsentry.ErrorKeys = []string{"error", "err"}
slogsentry.LogLevels = map[slog.Level]sentry.Level{...}
```

### Supported attributes

The following attributes are interpreted by `slogsentry.DefaultConverter`:

| Atribute name    | `slog.Kind`       | Underlying type |
| ---------------- | ----------------- | --------------- |
| "dist"           | string            |                 |
| "environment"    | string            |                 |
| "event_id"       | string            |                 |
| "platform"       | string            |                 |
| "release"        | string            |                 |
| "server_name"    | string            |                 |
| "tags"           | group (see below) |                 |
| "transaction"    | string            |                 |
| "user"           | group (see below) |                 |
| "error"          | any               | `error`         |
| "request"        | any               | `*http.Request` |
| "fingerprint"    | any               | `[]string`      |
| other attributes | *                 |                 |

Other attributes will be injected in `context` Sentry field.

Users and tags must be of type `slog.Group`. Eg:

```go
slog.Group("user",
    slog.String("id", "user-123"),
    slog.String("username", "samber"),
    slog.Time("created_at", time.Now()),
)
```

The Sentry agent is responsible for collecting `modules`.

### Example

```go
import (
    "github.com/getsentry/sentry-go"
    slogsentry "github.com/samber/slog-sentry/v2"
    "log/slog"
)

func main() {
    err := sentry.Init(sentry.ClientOptions{
        Dsn:           myDSN,
        EnableTracing: false,
    })
    if err != nil {
        log.Fatal(err)
    }

    defer sentry.Flush(2 * time.Second)

    logger := slog.New(slogsentry.Option{Level: slog.LevelDebug}.NewSentryHandler())
    logger = logger.
        With("environment", "dev").
        With("release", "v1.0.0")

    // log error
    logger.
        With("category", "sql").
        With("query.statement", "SELECT COUNT(*) FROM users;").
        With("query.duration", 1*time.Second).
        With("error", fmt.Errorf("could not count users")).
        Error("caramba!")

    // log user request
    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        With("request", httpRequest)
        With("status", 200).
        Info("received http request")
}
```

### Tracing

Import the samber/slog-otel library.

```go
import (
	slogsentry "github.com/samber/slog-sentry"
	slogotel "github.com/samber/slog-otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
	)
	tracer := tp.Tracer("hello/world")

	ctx, span := tracer.Start(context.Background(), "foo")
	defer span.End()

	span.AddEvent("bar")

	logger := slog.New(
		slogsentry.Option{
			// ...
			AttrFromContext: []func(ctx context.Context) []slog.Attr{
				slogotel.ExtractOtelAttrFromContext([]string{"tracing"}, "trace_id", "span_id"),
			},
		}.NewSentryHandler(),
	)

	logger.ErrorContext(ctx, "a message")
}
```

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/slog-sentry)
- Fix [open issues](https://github.com/samber/slog-sentry/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/slog-sentry)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
