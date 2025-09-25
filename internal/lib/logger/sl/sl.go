package sl

import (
	"log/slog"

)

// функция врзвращает атрибут пакета slog
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
