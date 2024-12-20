package slogsentry

import (
	"context"

	"log/slog"

	"github.com/getsentry/sentry-go"
	slogcommon "github.com/samber/slog-common"
)

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

func (o Option) NewSentryHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.AttrFromContext == nil {
		o.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
	}

	return &SentryHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*SentryHandler)(nil)

type SentryHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *SentryHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *SentryHandler) Handle(ctx context.Context, record slog.Record) error {
	hub := sentry.CurrentHub()
	if hubFromContext := sentry.GetHubFromContext(ctx); hubFromContext != nil {
		hub = hubFromContext
	} else if h.option.Hub != nil {
		hub = h.option.Hub
	}

	fromContext := slogcommon.ContextExtractor(ctx, h.option.AttrFromContext)
	event := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, append(h.attrs, fromContext...), h.groups, &record, hub)

	if h.option.BeforeSend != nil {
		event = h.option.BeforeSend(event)
	}

	hub.CaptureEvent(event)

	return nil
}

func (h *SentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SentryHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *SentryHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &SentryHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
