package masking_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/sodefrin/masking"
)

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
		logger.Info("test", masking.JSON("example", &Ex2{E: Ex{A: "a", B: "b", C: "c"}}, masking.Masking(Ex{}, "A", "B")))
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
