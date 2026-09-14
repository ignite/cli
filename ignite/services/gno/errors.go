package gno

import (
	"context"
	"log/slog"
	"strings"
	"sync"
)

// errorLog records ERROR level log records passing through a slog handler,
// so Serve can surface genesis failures after the node boots (the dev chain
// skips failing genesis transactions, so they would otherwise be silent).
type errorLog struct {
	mu   sync.Mutex
	errs []string
}

// record stores an error message.
func (e *errorLog) record(msg string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errs = append(e.errs, msg)
}

// errors returns the recorded error messages.
func (e *errorLog) errors() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.errs...)
}

// errorLogHandler is a slog handler that records ERROR records into an
// errorLog and delegates everything to the wrapped handler.
type errorLogHandler struct {
	log   *errorLog
	inner slog.Handler
}

// newErrorLogHandler wraps next with ERROR recording into log.
func newErrorLogHandler(log *errorLog, next slog.Handler) slog.Handler {
	return &errorLogHandler{log: log, inner: next}
}

// Enabled implements slog.Handler.
func (h *errorLogHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.inner.Enabled(ctx, l)
}

// Handle implements slog.Handler.
func (h *errorLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError {
		h.log.record(formatRecord(r))
	}
	return h.inner.Handle(ctx, r)
}

// formatRecord renders a record as "message key=value key=value".
func formatRecord(r slog.Record) string {
	var b strings.Builder
	b.WriteString(r.Message)
	r.Attrs(func(a slog.Attr) bool {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		b.WriteString(a.Value.String())
		return true
	})
	return b.String()
}

// WithAttrs implements slog.Handler.
func (h *errorLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &errorLogHandler{log: h.log, inner: h.inner.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (h *errorLogHandler) WithGroup(name string) slog.Handler {
	return &errorLogHandler{log: h.log, inner: h.inner.WithGroup(name)}
}
