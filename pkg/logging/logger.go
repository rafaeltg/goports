package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/rafaeltg/goports/pkg/cid"
)

type (
	handler struct {
		json slog.Handler
	}

	Configuration struct {
		slog.HandlerOptions
		fields []slog.Attr
	}

	ConfigurationOption func(*Configuration)
)

func NewLogger(opts ...ConfigurationOption) *slog.Logger {
	var cfg Configuration
	for _, opt := range opts {
		opt(&cfg)
	}

	var h slog.Handler

	h = slog.NewJSONHandler(
		os.Stdout,
		&cfg.HandlerOptions,
	)

	if len(cfg.fields) > 0 {
		h = h.WithAttrs(cfg.fields)
	}

	return slog.New(&handler{h})
}

func (h *handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.json.Enabled(ctx, lvl)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if cid, ok := cid.FromContext(ctx); ok {
		r.AddAttrs(slog.String("correlationId", cid))
	}

	return h.json.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.json.WithAttrs(attrs)
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h.json.WithGroup(name)
}

func WithSource(v bool) ConfigurationOption {
	return func(c *Configuration) {
		c.AddSource = v
	}
}

func WithLevel(l int) ConfigurationOption {
	return func(c *Configuration) {
		c.Level = slog.Level(l)
	}
}

func WithField(key string, value any) ConfigurationOption {
	return func(c *Configuration) {
		c.fields = append(c.fields, slog.Any(key, value))
	}
}
