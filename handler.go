package slogsentry

import (
	"context"

	"log/slog"

	"github.com/getsentry/sentry-go"
	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler
	// sentry hub (default: current hub)
	Hub *sentry.Hub

	// optional: customize Sentry event builder
	Converter Converter
	// optional: fetch attributes from context
	AttrFromContext []func(ctx context.Context) []slog.Attr

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
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
	return &SentryHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
