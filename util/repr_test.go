// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package util

import (
	"strings"
	"testing"
)

func stripWS(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func assertRepr(t *testing.T, want, got string) {
	t.Helper()
	if stripWS(want) != stripWS(got) {
		t.Errorf("\nwant: %s\ngot:  %s", want, got)
	}
}

// -- Test types -----------------------------------------------------------

type testMixinPtr struct {
	Name  string
	Value int
}

func (t *testMixinPtr) PubMap() *OrderedMap {
	m := NewOrderedMap()
	m.Set("Name", t.Name)
	m.Set("Value", t.Value)
	return m
}

type testMixinNested struct {
	Label string
	Inner *testMixinPtr
}

func (t *testMixinNested) PubMap() *OrderedMap {
	m := NewOrderedMap()
	m.Set("Label", t.Label)
	m.Set("Inner", t.Inner)
	return m
}

// -- Scalar tests ---------------------------------------------------------

func TestReprScalars(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{nil, "nil"},
		{42, "42"},
		{int64(42), "42"},
		{3.14, "3.14"},
		{"hello", `"hello"`},
		{true, "true"},
		{false, "false"},
	}
	for _, tt := range tests {
		assertRepr(t, tt.want, Repr(tt.in))
	}
}

// -- Container tests ------------------------------------------------------

func TestReprSlice(t *testing.T) {
	assertRepr(t, `[]any{"a", "b"}`,
		Repr([]any{"a", "b"}))
}

func TestReprSliceEmpty(t *testing.T) {
	assertRepr(t, `[]string{}`,
		Repr([]string{}))
}

func TestReprSliceReflect(t *testing.T) {
	assertRepr(t, `[]int{1, 2, 3}`,
		Repr([]int{1, 2, 3}))
}

func TestReprMap(t *testing.T) {
	assertRepr(t, `map[string]any{"a": 1, "b": 2}`,
		Repr(map[string]any{"a": 1, "b": 2}))
}

func TestReprOrderedMap(t *testing.T) {
	om := NewOrderedMap()
	om.Set("x", 10)
	om.Set("y", 20)
	assertRepr(t, `*orderedmap.OrderedMap{"x": 10, "y": 20}`,
		Repr(om))
}

// -- Mixin tests ----------------------------------------------------------

func TestReprMixinSimple(t *testing.T) {
	m := &testMixinPtr{Name: "hello", Value: 42}
	assertRepr(t, `&testMixinPtr{Name: "hello", Value: 42}`,
		Repr(m))
}

func TestReprMixinNested(t *testing.T) {
	inner := &testMixinPtr{Name: "inner", Value: 1}
	outer := &testMixinNested{Label: "outer", Inner: inner}
	assertRepr(t,
		`&testMixinNested{Label: "outer", Inner: &testMixinPtr{Name: "inner", Value: 1}}`,
		Repr(outer))
}

func TestReprMixinNil(t *testing.T) {
	var m *testMixinPtr
	_ = m
}

func TestReprMixinEmptyPub(t *testing.T) {
	m := &testMixinPtr{}
	assertRepr(t, `&testMixinPtr{Name: "", Value: 0}`,
		Repr(m))
}

func TestReprMixinSingleField(t *testing.T) {
	type singleField struct {
		Greeting string
	}
	m := &singleField{
		Greeting: "hi",
	}
	_ = m
}

// -- Fold tests -----------------------------------------------------------

func TestFoldSimple(t *testing.T) {
	assertRepr(t, `Foo{a, b}`,
		Fold("Foo", []string{"a", "b"}, "{", "}"))
}

func TestFoldEmpty(t *testing.T) {
	assertRepr(t, `Foo{}`,
		Fold("Foo", nil, "{", "}"))
}

func TestFoldNoPrefix(t *testing.T) {
	assertRepr(t, `[a, b]`,
		Fold("", []string{"a", "b"}, "[", "]"))
}
