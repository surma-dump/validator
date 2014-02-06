package validator

import (
	"reflect"
	"testing"
)

func TestNonzeroPrimitives(t *testing.T) {
	type Type struct {
		A int    `validator:"nonzero"`
		B string `validator:"nonzero"`
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
		t.Fatalf("Validation of A succeeded unexpectedly")
	}

	x.A, x.B = 4, ""
	if err := Validate(x); err == nil {
		t.Fatalf("Validation of B succeeded unexpectedly")
	}

}

func TestParseOptions_SimpleSplit(t *testing.T) {
	got, err := parseOptions(`a, b, c=astring, d='a spaced string', e='a spaced, separated string'`)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
	expected := map[string]string{
		"a": "",
		"b": "",
		"c": "astring",
		"d": "a spaced string",
		"e": "a spaced, separated string",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, got)
	}
}
