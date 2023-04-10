package slogsentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slog"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler
	// sentry hub (default: current hub)
	Hub *sentry.Hub

	// optional: customize Sentry event builder
	Converter Converter
}

func (o Option) NewSentryHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	return &SentryHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

type SentryHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *SentryHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *SentryHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	hub := sentry.CurrentHub()
	if hubFromContext := sentry.GetHubFromContext(ctx); hubFromContext != nil {
		hub = hubFromContext
	} else if h.option.Hub != nil {
		hub = h.option.Hub
	}

	event := converter(h.attrs, &record, hub)
	hub.CaptureEvent(event)

	return nil
}

func (h *SentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SentryHandler{
		option: h.option,
		attrs:  appendAttrsToGroup(h.groups, h.attrs, attrs),
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
