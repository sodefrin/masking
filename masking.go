package masking

import (
	"log/slog"

	"github.com/sodefrin/masking/internal/json"
)

type Option json.MarshalOption

func WithFieldsFilter(typ any, fields ...string) Option {
	return Option(json.WithStructFieldFilter(typ, fields...))
}

func SetMaskedString(s string) Option {
	return Option(json.WithFilterString(s))
}

func WithTagFilter(tagKey, tagValue string) Option {
	return Option(json.WithTagFilter(tagKey, tagValue))
}

func JSON(k string, v any, opts ...Option) slog.Attr {
	jsonOpts := make([]json.MarshalOption, len(opts))
	for i, opt := range opts {
		jsonOpts[i] = json.MarshalOption(opt)
	}
	return slog.Any(k, customJSON{v, jsonOpts})
}

type customJSON struct {
	value any
	ops   []json.MarshalOption
}

func (t customJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value, t.ops...)
}
