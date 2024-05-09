package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rntrp/mailheap/internal/config"
)

func Logger() *slog.Logger {
	var leveler slog.Leveler
	switch strings.ToUpper(config.GetLogLevel()) {
	case "DEBUG":
		leveler = slog.LevelDebug
	case "WARN":
		leveler = slog.LevelWarn
	case "ERROR":
		leveler = slog.LevelError
	default:
		leveler = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     leveler,
	}
	w := os.Stdout
	switch strings.ToUpper(config.GetLogFormat()) {
	case "JSON":
		return slog.New(slog.NewJSONHandler(w, opts))
	case "TEXT":
		return slog.New(slog.NewTextHandler(w, opts))
	case "SIMPLE":
		return slog.New(newSimpleHandler(w, opts))
	default:
		panic("Unknown logger" + config.GetLogFormat())
	}
}

type simpleHandler struct {
	h  slog.Handler
	w  io.Writer
	mu *sync.Mutex
}

func newSimpleHandler(w io.Writer, opts *slog.HandlerOptions) *simpleHandler {
	return &simpleHandler{
		h: slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: nil,
		}),
		w:  w,
		mu: &sync.Mutex{},
	}
}

func (h *simpleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *simpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &simpleHandler{h: h.h.WithAttrs(attrs), w: h.w, mu: h.mu}
}

func (h *simpleHandler) WithGroup(name string) slog.Handler {
	return &simpleHandler{h: h.h.WithGroup(name), w: h.w, mu: h.mu}
}

var linefeed = []byte("\n")

func (h *simpleHandler) Handle(ctx context.Context, r slog.Record) error {
	t := r.Time.Format(time.RFC3339)
	v := []string{t, r.Message}
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			v = append(v, a.Value.String())
			return true
		})
	}
	line := []byte(strings.Join(v, " "))
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := h.w.Write(line); err != nil {
		return fmt.Errorf("simpleHandler line: %w", err)
	} else if _, err := h.w.Write(linefeed); err != nil {
		return fmt.Errorf("simpleHandler linefeed: %w", err)
	}
	return nil
}
