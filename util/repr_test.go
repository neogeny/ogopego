package util

import (
	"testing"

	"github.com/iancoleman/orderedmap"
)

type ReprTestStruct struct {
	Name  string
	Value int
	Ok    bool
}

type ReprNested struct {
	Label string
	Inner *ReprTestStruct
}

func TestReprPrimitives(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{nil, "nil"},
		{true, "true"},
		{false, "false"},
		{42, "42"},
		{int8(8), "8"},
		{uint(1), "1"},
		{3.14, "3.14"},
		{"hello", `"hello"`},
	}
	for _, tt := range tests {
		got := Repr(tt.input)
		if got != tt.want {
			t.Errorf("Repr(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestReprStruct(t *testing.T) {
	s := &ReprTestStruct{Name: "foo", Value: 7, Ok: true}
	got := Repr(s)
	want := `&util.ReprTestStruct{Name: "foo", Value: 7, Ok: true}`
	if got != want {
		t.Errorf("Repr(%v) = %q, want %q", s, got, want)
	}
}

func TestReprSlice(t *testing.T) {
	got := Repr([]any{1, "two", true})
	want := `[]any{1, "two", true}`
	if got != want {
		t.Errorf("Repr() = %q, want %q", got, want)
	}
}

func TestReprFlatMap(t *testing.T) {
	got := Repr(map[string]any{"a": 1, "b": 2})
	want := `map[string]any{a: 1, b: 2}`
	if got != want {
		t.Errorf("Repr(map) = %q, want %q", got, want)
	}
}

func TestReprMapWithClass(t *testing.T) {
	m := map[string]any{
		"__class__": "MyType",
		"name":      "foo",
		"value":     42,
	}
	got := Repr(m)
	want := `MyType{name: "foo", value: 42}`
	if got != want {
		t.Errorf("Repr(__class__ map) = %q, want %q", got, want)
	}
}

func TestReprNestedStruct(t *testing.T) {
	s := &ReprNested{
		Label: "top",
		Inner: &ReprTestStruct{Name: "inner", Value: 1, Ok: false},
	}
	got := Repr(s)
	want := "&util.ReprNested{\n  Label: \"top\",\n  Inner: &util.ReprTestStruct{Name: \"inner\", Value: 1, Ok: false},\n}"
	if got != want {
		t.Errorf("Repr(%v) = %q, want %q", s, got, want)
	}
}

func TestReprTypedSlice(t *testing.T) {
	type item struct{ Name string }
	got := Repr([]*item{{Name: "a"}, {Name: "b"}})
	want := `[]*util.item{&util.item{Name: "a"}, &util.item{Name: "b"}}`
	if got != want {
		t.Errorf("Repr([]*item) = %q, want %q", got, want)
	}
}

func TestReprTypedStringSlice(t *testing.T) {
	got := Repr([]string{"x", "y"})
	want := `[]string{"x", "y"}`
	if got != want {
		t.Errorf("Repr([]string) = %q, want %q", got, want)
	}
}

func TestReprOrderedMap(t *testing.T) {
	om := orderedmap.New()
	om.Set("count", 42)
	om.Set("hello", "world")
	got := Repr(om)
	want := `map[string]any{count: 42, hello: "world"}`
	if got != want {
		t.Errorf("Repr(OrderedMap) = %q, want %q", got, want)
	}
}

func TestReprEmptyOrderedMap(t *testing.T) {
	om := orderedmap.New()
	got := Repr(om)
	want := "map[string]any{}"
	if got != want {
		t.Errorf("Repr(empty OrderedMap) = %q, want %q", got, want)
	}
}

func TestReprPointerToInt(t *testing.T) {
	n := 42
	got := Repr(&n)
	want := "&42"
	if got != want {
		t.Errorf("Repr(&int) = %q, want %q", got, want)
	}
}
