package masking

import (
	"log/slog"
)

func JSON(k string, v any, opts ...Option) slog.Attr {
	return slog.Any(k, customJSON{v, opts})
}

type customJSON struct {
	value any
	ops   []Option
}

func (t customJSON) MarshalJSON() ([]byte, error) {
	return MarshalJSON(t.value, t.ops...)
}
