package masking_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/sodefrin/masking"
	"github.com/sodefrin/masking/internal/json"
)

type Ex1 struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

type Ex2 struct {
	E Ex1 `json:"e"`
}

func TestMasking(t *testing.T) {
	ret, err := json.Marshal(&Ex1{A: "a", B: "b", C: "c"}, json.WithStructFieldFilter(Ex1{}, "a", "b"), json.WithFilterString("xxx"))
	if err != nil {
		t.Error(err)
	}
	if want, have := `{"a":"xxx","b":"xxx","c":"c"}`, string(ret); want != have {
		t.Errorf("unexpected marshal: want %v have %v", want, have)
	}
	ret, err = json.Marshal(&Ex2{E: Ex1{A: "a", B: "b", C: "c"}}, json.WithStructFieldFilter(Ex1{}, "a", "b"), json.WithFilterString("xxx"))
	if err != nil {
		t.Error(err)
	}
	if want, have := `{"e":{"a":"xxx","b":"xxx","c":"c"}}`, string(ret); want != have {
		t.Errorf("unexpected marshal: want %v have %v", want, have)
	}
}

func TestJSON(t *testing.T) {
	r, w := io.Pipe()

	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.String("time", "dummy")
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{ReplaceAttr: replace}))
	go func() {
		logger.Info("test", masking.JSON("example", &Ex2{E: Ex1{A: "a", B: "b", C: "c"}}, masking.IgnoreFields(Ex1{}, "a", "b"), masking.SetMaskedString("xxx")))
		w.Close()
	}()

	ret, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := `{"time":"dummy","level":"INFO","msg":"test","example":{"e":{"a":"xxx","b":"xxx","c":"c"}}}`, string(ret)[:len(string(ret))-1]; want != have {
		t.Fatalf("unexpected marshal: want %v have %v", want, have)
	}
}
