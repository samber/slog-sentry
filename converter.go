package slogsentry

import (
	"net/http"
	"reflect"

	"log/slog"

	"github.com/getsentry/sentry-go"
	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"
var ContextKey = "extra"
var ErrorKeys = []string{"error", "err"}

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record, hub *sentry.Hub) *sentry.Event

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record, hub *sentry.Hub) *sentry.Event {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	attrs = slogcommon.ReplaceError(attrs, ErrorKeys...)
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)

	// handler formatter
	event := sentry.NewEvent()
	event.Timestamp = record.Time.UTC()
	event.Level = LogLevels[record.Level]
	event.Message = record.Message
	event.Logger = name

	for i := range attrs {
		attrToSentryEvent(attrs[i], event)
	}

	return event
}

func attrToSentryEvent(attr slog.Attr, event *sentry.Event) {
	k := attr.Key
	v := attr.Value
	kind := attr.Value.Kind()

	for _, errorKey := range ErrorKeys {
		if attr.Key == errorKey {
			if err, ok := attr.Value.Any().(error); ok {
				event.Exception = buildExceptions(err)
			} else {
				event.User.Data[errorKey] = slogcommon.AnyValueToString(v)
			}
		}
	}

	if k == "dist" && kind == slog.KindString {
		event.Dist = v.String()
	} else if k == "environment" && kind == slog.KindString {
		event.Environment = v.String()
	} else if k == "event_id" && kind == slog.KindString {
		event.EventID = sentry.EventID(v.String())
	} else if k == "platform" && kind == slog.KindString {
		event.Platform = v.String()
	} else if k == "release" && kind == slog.KindString {
		event.Release = v.String()
	} else if k == "server_name" && kind == slog.KindString {
		event.ServerName = v.String()
	} else if attr.Key == "tags" && kind == slog.KindGroup {
		event.Tags = slogcommon.AttrsToString(v.Group()...)
	} else if attr.Key == "transaction" && kind == slog.KindGroup {
		event.Transaction = v.String()
	} else if attr.Key == "user" && kind == slog.KindGroup {
		data := slogcommon.AttrsToString(v.Group()...)

		if id, ok := data["id"]; ok {
			event.User.ID = id
			delete(data, "id")
		} else if email, ok := data["email"]; ok {
			event.User.Email = email
			delete(data, "email")
		} else if ipAddress, ok := data["ip_address"]; ok {
			event.User.IPAddress = ipAddress
			delete(data, "ip_address")
		} else if username, ok := data["username"]; ok {
			event.User.Username = username
			delete(data, "username")
		} else if name, ok := data["name"]; ok {
			event.User.Name = name
			delete(data, "name")
		} else if segment, ok := data["segment"]; ok {
			event.User.Segment = segment
			delete(data, "segment")
		}

		event.User.Data = data
	} else if attr.Key == "request" && kind == slog.KindAny {
		if req, ok := attr.Value.Any().(http.Request); ok {
			event.Request = sentry.NewRequest(&req)
		} else if req, ok := attr.Value.Any().(*http.Request); ok {
			event.Request = sentry.NewRequest(req)
		} else {
			event.User.Data["request"] = slogcommon.AnyValueToString(v)
		}
	} else if kind == slog.KindGroup {
		event.Contexts[attr.Key] = slogcommon.AttrsToMap(attr.Value.Group()...)
	} else {
		// "context" should not be added to underlying context layers (see slog.KindGroup case).
		if _, ok := event.Contexts[ContextKey]; !ok {
			event.Contexts[ContextKey] = make(map[string]any, 0)
		}
		event.Contexts[ContextKey][attr.Key] = attr.Value.Any()
	}
}

func buildExceptions(err error) []sentry.Exception {
	exceptions := []sentry.Exception{}

	for i := 0; i < 10 && err != nil; i++ {
		exceptions = append(exceptions, sentry.Exception{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: sentry.ExtractStacktrace(err), // @TODO: use record.pc
		})

		switch previous := err.(type) {
		case interface{ Unwrap() error }:
			err = previous.Unwrap()
		case interface{ Cause() error }:
			err = previous.Cause()
		default:
			err = nil
		}
	}

	return exceptions
}
