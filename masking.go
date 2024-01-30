package masking

import (
	"reflect"
	"sync"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var defaultMaskingStr = ""

func MarshalJSON(in any, opts ...Option) ([]byte, error) {
	cfg := jsoniter.Config{}.Froze()

	for _, o := range opts {
		for _, e := range o() {
			cfg.RegisterExtension(e)
		}
	}

	return cfg.Marshal(in)
}

func SetMaskedStr(s string) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	defaultMaskingStr = s
}

type Option func() []jsoniter.Extension

func IgnoreFields(typ any, fields ...string) Option {
	return func() []jsoniter.Extension {
		es := []jsoniter.Extension{}
		for _, f := range fields {
			if e := newIgnoreFieldsExtension(typ, f, defaultMaskingStr); e != nil {
				es = append(es, e)
			}
		}
		return es
	}
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}

func newIgnoreFieldsExtension(typ interface{}, field, maskedStr string) *ignoreFieldsExtension {
	t := reflect.TypeOf(typ)
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	return &ignoreFieldsExtension{
		typ:       t.String(),
		field:     field,
		maskedStr: maskedStr,
	}
}

type ignoreFieldsExtension struct {
	jsoniter.DummyExtension
	typ       string
	field     string
	maskedStr string
}

func (m *ignoreFieldsExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	if structDescriptor.Type.String() != m.typ {
		return
	}

	binding := structDescriptor.GetField(m.field)
	binding.Encoder = &funcEncoder{
		fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			stream.WriteString(m.maskedStr)
		},
	}
}
