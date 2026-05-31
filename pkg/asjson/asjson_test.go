// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package asjson

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestAsJSONPrimitives(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{nil, "null"},
		{true, "true"},
		{false, "false"},
		{42, "42"},
		{int64(42), "42"},
		{3.14, "3.14"},
		{"hello", `"hello"`},
	}
	for _, tt := range tests {
		got := AsJSONStr(tt.in)
		assert.Equal(t, tt.want, got, "AsJSONStr(%v) = %s, want %s", tt.in, got, tt.want)
	}
}

func TestAsJSONSlice(t *testing.T) {
	v := []any{1, "two", true}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	assert.Equal(t, `[1,"two",true]`, string(b))
}

func TestAsJSONMap(t *testing.T) {
	v := map[string]any{"a": 1, "b": "hello"}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	assert.Equal[any](t, float64(1), out["a"], "expected 1, got %v", out["a"])
}

func TestAsJSONNestedSlice(t *testing.T) {
	v := [][]any{{1, 2}, {3, 4}}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	assert.Equal(t, `[[1,2],[3,4]]`, string(b))
}

type testStruct struct {
	Name   string
	Value  int
	hidden string
}

func TestAsJSONStruct(t *testing.T) {
	v := testStruct{Name: "foo", Value: 42, hidden: "secret"}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	assert.Equal[any](t, "foo", out["name"], "expected foo, got %v", out["Name"])
	assert.Equal[any](t, float64(42), out["value"], "expected 42, got %v", out["Value"])
	_, ok := out["hidden"]
	assert.False(t, ok, "unexported field 'hidden' should not appear")
}

type jsonStruct struct {
	Name string
}

func (s *jsonStruct) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"__class__": "jsonStruct", "name": s.Name})
}

func TestAsJSONJSONMarshaler(t *testing.T) {
	v := &jsonStruct{Name: "test"}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	assert.Equal[any](t, "asjson.jsonStruct", out["__class__"], "expected __class__ asjson.jsonStruct, got %v", out["__class__"])
	assert.Equal[any](t, "test", out["name"], "expected test, got %v", out["name"])
}

func TestAsJSONCycle(t *testing.T) {
	type node struct {
		Name string
		Next *node
	}
	a := &node{Name: "a"}
	b := &node{Name: "b"}
	a.Next = b
	b.Next = a

	result := AsJSON(a)
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected map, got %T", result)
	assert.Equal[any](t, "a", m["name"], "expected 'a', got %v", m["Name"])
	next, ok := m["next"].(map[string]any)
	assert.True(t, ok, "expected map for Next, got %T", m["Next"])
	assert.Equal[any](t, "b", next["name"], "expected 'b', got %v", next["Name"])
	cycle, ok := next["next"].(string)
	assert.True(t, ok, "expected string cycle marker, got %T", next["Next"])
	assert.True(t, len(cycle) >= 3, "expected meaningful cycle marker, got %q", cycle)
}

func TestAsJSONs(t *testing.T) {
	v := map[string]any{"x": 1}
	s := AsJSONStr(v)
	want := "{\n  \"x\": 1\n}"
	assert.Equal(t, want, s, "got %q, want %q", s, want)
}

func TestAsJSONChanFunc(t *testing.T) {
	ch := make(chan int)
	fn := func() {}
	v := []any{ch, fn}
	b, err := json.Marshal(AsJSON(v))
	assert.NoError(t, err)
	assert.Equal(t, `[null,null]`, string(b))
}
