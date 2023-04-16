
# slog: Sentry handler

[![tag](https://img.shields.io/github/tag/samber/slog-sentry.svg)](https://github.com/samber/slog-sentry/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20.1-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/slog-sentry?status.svg)](https://pkg.go.dev/github.com/samber/slog-sentry)
![Build Status](https://github.com/samber/slog-sentry/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/slog-sentry)](https://goreportcard.com/report/github.com/samber/slog-sentry)
[![Coverage](https://img.shields.io/codecov/c/github/samber/slog-sentry)](https://codecov.io/gh/samber/slog-sentry)
[![Contributors](https://img.shields.io/github/contributors/samber/slog-sentry)](https://github.com/samber/slog-sentry/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/slog-sentry)](./LICENSE)

A [Sentry](https://sentry.io) Handler for [slog](https://pkg.go.dev/golang.org/x/exp/slog) Go library.

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): workflows of `slog` handlers (pipeline, fanout, ...)
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`

## 🚀 Install

```sh
go get github.com/samber/slog-sentry
```

**Compatibility**: go >= 1.20.1

This library is v0 and follows SemVer strictly. On `slog` final release (go 1.21), this library will go v1.

No breaking changes will be made to exported APIs before v1.0.0.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/slog-sentry](https://pkg.go.dev/github.com/samber/slog-sentry)

### Handler options

```go
type Option struct {
    // log level (default: debug)
	Level     slog.Leveler
    // sentry hub (default: current hub)
	Hub       *sentry.Hub

    // optional: customize Sentry event builder
	Converter Converter
}
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
| other attributes | *                 |                 |

Other attributes will be injected in `extra` Sentry field.

Users and tags must be of type `slog.Group`. Eg:

```go
slog.Group("user",
    slog.String("id", "user-123"),
    slog.String("username", "samber"),
    slog.Time("created_at", time.Now()),
)
```

### Example

```go
import (
	"github.com/getsentry/sentry-go"
	slogsentry "github.com/samber/slog-sentry"
	"golang.org/x/exp/slog"
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
