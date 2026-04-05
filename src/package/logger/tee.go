package logger

import (
	"context"
	"errors"
	"log/slog"
)

// teeHandler duplicates log records to two handlers (e.g. JSON file + text console).
type teeHandler struct {
	a, b slog.Handler
}

func newTeeHandler(a, b slog.Handler) slog.Handler {
	return &teeHandler{a: a, b: b}
}

func (t *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.a.Enabled(ctx, level) || t.b.Enabled(ctx, level)
}

func (t *teeHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	if t.a.Enabled(ctx, r.Level) {
		errs = append(errs, t.a.Handle(ctx, r.Clone()))
	}
	if t.b.Enabled(ctx, r.Level) {
		errs = append(errs, t.b.Handle(ctx, r.Clone()))
	}
	return errors.Join(errs...)
}

func (t *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &teeHandler{a: t.a.WithAttrs(attrs), b: t.b.WithAttrs(attrs)}
}

func (t *teeHandler) WithGroup(name string) slog.Handler {
	return &teeHandler{a: t.a.WithGroup(name), b: t.b.WithGroup(name)}
}
