package slogdiscard

import (
	"context"
	"log/slog"
)

// конструктор логгера
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// обработчик
type DiscardHandler struct{}

// конструктор обработчика
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil // игнорируем записи журнала
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false // false, так как игнорируем записи журнала
}