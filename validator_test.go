package validator

import (
	"testing"
)

func TestNonzeroPrimitives(t *testing.T) {
	type Type struct {
		A int
		B string
	}

	x := Type{
		A: 4,
		B: "hai",
	}

	if err := Validate(x); err != nil {
		t.Fatalf("Validation failed: %s", err)
	}

	x.A = 0
	if err := Validate(x); err == nil {
		t.Fatalf("Validation succeeded unexpectedly")
	}

	x.A, x.B = 4, ""
	if err := Validate(x); err == nil {
		t.Fatalf("Validation succeeded unexpectedly")
	}

}
