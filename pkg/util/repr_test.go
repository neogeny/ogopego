package util

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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
		assert.Equal(t, tt.want, got, "Repr(%v)", tt.input)
	}
}

func TestReprStruct(t *testing.T) {
	s := &ReprTestStruct{Name: "foo", Value: 7, Ok: true}
	got := Repr(s)
	want := `&util.ReprTestStruct{Name: "foo", Value: 7, Ok: true}`
	assert.Equal(t, want, got, "Repr(%v)", s)
}

func TestReprSlice(t *testing.T) {
	got := Repr([]any{1, "two", true})
	want := `[]any{1, "two", true}`
	assert.Equal(t, want, got)
}

func TestReprFlatMap(t *testing.T) {
	got := Repr(map[string]any{"a": 1, "b": 2})
	want := `map[string]any{a: 1, b: 2}`
	assert.Equal(t, want, got)
}

func TestReprMapWithClass(t *testing.T) {
	m := map[string]any{
		"__class__": "MyType",
		"name":      "foo",
		"value":     42,
	}
	got := Repr(m)
	want := `MyType{name: "foo", value: 42}`
	assert.Equal(t, want, got)
}

func TestReprNestedStruct(t *testing.T) {
	s := &ReprNested{
		Label: "top",
		Inner: &ReprTestStruct{Name: "inner", Value: 1, Ok: false},
	}
	got := Repr(s)
	want := "&util.ReprNested{\n  Label: \"top\",\n  Inner: &util.ReprTestStruct{Name: \"inner\", Value: 1, Ok: false},\n}"
	assert.Equal(t, want, got, "Repr(%v)", s)
}

func TestReprTypedSlice(t *testing.T) {
	type item struct{ Name string }
	got := Repr([]*item{{Name: "a"}, {Name: "b"}})
	want := `[]*util.item{&util.item{Name: "a"}, &util.item{Name: "b"}}`
	assert.Equal(t, want, got)
}

func TestReprTypedStringSlice(t *testing.T) {
	got := Repr([]string{"x", "y"})
	want := `[]string{"x", "y"}`
	assert.Equal(t, want, got)
}

func TestReprOrderedMap(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("count", 42)
	om.Set("hello", "world")
	got := Repr(om)
	want := `map[string]any{count: 42, hello: "world"}`
	assert.Equal(t, want, got)
}

func TestReprEmptyOrderedMap(t *testing.T) {
	om := orderedmap.New[string, any]()
	got := Repr(om)
	want := "map[string]any{}"
	assert.Equal(t, want, got)
}

func TestReprPointerToInt(t *testing.T) {
	n := 42
	got := Repr(&n)
	want := "&42"
	assert.Equal(t, want, got)
}
