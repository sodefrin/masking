## masking

masking provides `cmpopts.IgnoreFields` like interface to mask JSON field for slog.

```Go
type Ex1 struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c" secret:"pii"`
}

logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{ReplaceAttr: replace}))
logger.Info(
	"test",
	masking.JSON(
		"json", Ex1{A: "a", B: "b", C: "c"},
		masking.IgnoreFields(Ex1{}, "a", "b"),
		masking.WithTagFilter("secret", "pii"),
		masking.SetMaskedString("xxx"),
	),
)
// {"level":"INFO","msg":"test","json"{"a":"xxx","b":"xxx","c":"xxx"}}`
```
