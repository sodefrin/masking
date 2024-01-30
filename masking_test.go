package masking_test

import (
	"testing"

	"github.com/sodefrin/masking"
)

type Ex struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

type Ex2 struct {
	E Ex `json:"e"`
}

func TestMasking(t *testing.T) {
	masking.SetMaskedStr("xxx")
	ret, err := masking.MarshalJSON(&Ex{A: "a", B: "b", C: "c"}, masking.Masking(Ex{}, "A", "B"))
	if err != nil {
		t.Fatal(err)
	}
	if want, have := `{"a":"xxx","b":"xxx","c":"c"}`, string(ret); want != have {
		t.Fatalf("unexpected marshal: want %v have %v", want, have)
	}
	ret, err = masking.MarshalJSON(&Ex2{E: Ex{A: "a", B: "b", C: "c"}}, masking.Masking(Ex{}, "A", "B"))
	if err != nil {
		t.Fatal(err)
	}
	if want, have := `{"e":{"a":"xxx","b":"xxx","c":"c"}}`, string(ret); want != have {
		t.Fatalf("unexpected marshal: want %v have %v", want, have)
	}
}
