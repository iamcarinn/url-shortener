package sl

import "log/slog"

// Возвращает атрибут для ошибки, кот. можно использовать в slog
func Err(err error) slog.Attr {
	return slog.Attr {
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}