package slogsentry

import (
	"encoding"
	"fmt"
	"net/http"

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
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)
	// Removes err attribute from `attrs`
	attrs, err := slogcommon.ExtractError(attrs, ErrorKeys...)

	// handler formatter
	event := sentry.NewEvent()
	event.Timestamp = record.Time.UTC()
	event.Level = LogLevels[record.Level]
	event.Message = record.Message
	event.Logger = name
	event.SetException(err, 10)

	for i := range attrs {
		attrToSentryEvent(attrs[i], event)
	}

	return event
}

func attrToSentryEvent(attr slog.Attr, event *sentry.Event) {
	k := attr.Key
	v := attr.Value
	kind := v.Kind()

	switch {
	case k == "dist" && kind == slog.KindString:
		event.Dist = v.String()
	case k == "environment" && kind == slog.KindString:
		event.Environment = v.String()
	case k == "event_id" && kind == slog.KindString:
		event.EventID = sentry.EventID(v.String())
	case k == "platform" && kind == slog.KindString:
		event.Platform = v.String()
	case k == "release" && kind == slog.KindString:
		event.Release = v.String()
	case k == "server_name" && kind == slog.KindString:
		event.ServerName = v.String()
	case k == "tags" && kind == slog.KindGroup:
		event.Tags = slogcommon.AttrsToString(v.Group()...)
	case k == "transaction" && kind == slog.KindString:
		event.Transaction = v.String()
	case k == "user" && kind == slog.KindGroup:
		data := slogcommon.AttrsToString(v.Group()...)

		if id, ok := data["id"]; ok {
			event.User.ID = id
			delete(data, "id")
		}
		if email, ok := data["email"]; ok {
			event.User.Email = email
			delete(data, "email")
		}
		if ipAddress, ok := data["ip_address"]; ok {
			event.User.IPAddress = ipAddress
			delete(data, "ip_address")
		}
		if username, ok := data["username"]; ok {
			event.User.Username = username
			delete(data, "username")
		}
		if name, ok := data["name"]; ok {
			event.User.Name = name
			delete(data, "name")
		}

		event.User.Data = data
	case k == "request" && kind == slog.KindAny:
		if req, ok := v.Any().(http.Request); ok {
			event.Request = sentry.NewRequest(&req)
		} else if req, ok := v.Any().(*http.Request); ok {
			event.Request = sentry.NewRequest(req)
		} else {
			if tm, ok := v.Any().(encoding.TextMarshaler); ok {
				data, err := tm.MarshalText()
				if err == nil {
					event.User.Data["request"] = string(data)
				} else {
					event.User.Data["request"] = fmt.Sprintf("%v", v.Any())
				}
			}
		}
	case k == "fingerprint" && kind == slog.KindAny:
		if fingerprint, ok := v.Any().([]string); ok {
			event.Fingerprint = fingerprint
		}
	case kind == slog.KindGroup:
		event.Contexts[k] = slogcommon.AttrsToMap(v.Group()...)
	default:
		// "context" should not be added to underlying context layers (see slog.KindGroup case).
		if _, ok := event.Contexts[ContextKey]; !ok {
			event.Contexts[ContextKey] = make(map[string]any, 0)
		}
		event.Contexts[ContextKey][k] = v.Any()
	}
}
