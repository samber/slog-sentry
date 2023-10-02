package slogsentry

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"log/slog"

	"github.com/getsentry/sentry-go"
)

type Converter func(loggerAttr []slog.Attr, record *slog.Record, hub *sentry.Hub) *sentry.Event

func DefaultConverter(loggerAttr []slog.Attr, record *slog.Record, hub *sentry.Hub) *sentry.Event {
	event := sentry.NewEvent()

	event.Timestamp = record.Time.UTC()
	event.Level = levelMap[record.Level]
	event.Message = record.Message
	event.Logger = "samber/slog-sentry"

	for i := range loggerAttr {
		attrToSentryEvent(loggerAttr[i], event)
	}

	record.Attrs(func(attr slog.Attr) bool {
		attrToSentryEvent(attr, event)
		return true
	})

	return event
}

func attrToSentryEvent(attr slog.Attr, event *sentry.Event) {
	k := attr.Key
	v := attr.Value
	kind := attr.Value.Kind()

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
		event.Tags = attrToStringMap(v.Group())
	} else if attr.Key == "transaction" && kind == slog.KindGroup {
		event.Transaction = v.String()
	} else if attr.Key == "user" && kind == slog.KindGroup {
		data := attrToStringMap(v.Group())

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
	} else if attr.Key == "error" && kind == slog.KindAny {
		if err, ok := attr.Value.Any().(error); ok {
			event.Exception = buildExceptions(err)
		} else {
			event.User.Data["error"] = anyValueToString(v)
		}
	} else if attr.Key == "request" && kind == slog.KindAny {
		if req, ok := attr.Value.Any().(http.Request); ok {
			event.Request = sentry.NewRequest(&req)
		} else if req, ok := attr.Value.Any().(*http.Request); ok {
			event.Request = sentry.NewRequest(req)
		} else {
			event.User.Data["request"] = anyValueToString(v)
		}
	} else if kind == slog.KindGroup {
		event.Contexts[attr.Key] = attrToMap(attr.Value.Group())
	} else {
		// "context" should not be added to underlying context layers (see slog.KindGroup case).
		if _, ok := event.Contexts["context"]; !ok {
			event.Contexts["context"] = map[string]any{}
		}
		event.Contexts["context"][attr.Key] = attr.Value.Any()
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

func attrToMap(attrs []slog.Attr) map[string]any {
	output := map[string]any{}
	for i := range attrs {
		attr := attrs[i]
		k := attr.Key
		v := attr.Value
		kind := attr.Value.Kind()

		switch kind {
		case slog.KindAny:
			output[k] = anyValueToString(v)
		case slog.KindLogValuer:
			output[k] = anyValueToString(v)
		case slog.KindGroup:
			output[k] = attrToMap(v.Group())
		case slog.KindInt64:
			output[k] = v.Int64()
		case slog.KindUint64:
			output[k] = v.Uint64()
		case slog.KindFloat64:
			output[k] = v.Float64()
		case slog.KindString:
			output[k] = v.String()
		case slog.KindBool:
			output[k] = v.Bool()
		case slog.KindDuration:
			output[k] = v.Duration()
		case slog.KindTime:
			output[k] = v.Time().UTC()
		default:
			output[k] = anyValueToString(v)
		}
	}
	return output
}

func attrToStringMap(attrs []slog.Attr) map[string]string {
	output := map[string]string{}
	for i := range attrs {
		attr := attrs[i]
		k, v := attr.Key, attr.Value
		output[k] = valueToString(v)
	}
	return output
}

func valueToString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindAny:
		return anyValueToString(v)
	case slog.KindLogValuer:
		return anyValueToString(v)
	case slog.KindGroup:
		return fmt.Sprint(v)
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().UTC().String()
	default:
		return anyValueToString(v)
	}
}

func anyValueToString(v slog.Value) string {
	if tm, ok := v.Any().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return ""
		}

		return string(data)
	}

	return fmt.Sprintf("%+v", v.Any())
}
