// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package json

import (
	"encoding/json"
	"testing"
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
		if got != tt.want {
			t.Errorf("AsJSONStr(%v) = %s, want %s", tt.in, got, tt.want)
		}
	}
}

func TestAsJSONSlice(t *testing.T) {
	v := []any{1, "two", true}
	b, err := json.Marshal(AsJSON(v))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `[1,"two",true]` {
		t.Errorf("got %s", string(b))
	}
}

func TestAsJSONMap(t *testing.T) {
	v := map[string]any{"a": 1, "b": "hello"}
	b, err := json.Marshal(AsJSON(v))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["a"] != float64(1) {
		t.Errorf("expected 1, got %v", out["a"])
	}
}

func TestAsJSONNestedSlice(t *testing.T) {
	v := [][]any{{1, 2}, {3, 4}}
	b, err := json.Marshal(AsJSON(v))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `[[1,2],[3,4]]` {
		t.Errorf("got %s", string(b))
	}
}

type testStruct struct {
	Name   string
	Value  int
	hidden string
}

func TestAsJSONStruct(t *testing.T) {
	v := testStruct{Name: "foo", Value: 42, hidden: "secret"}
	b, err := json.Marshal(AsJSON(v))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["name"] != "foo" {
		t.Errorf("expected foo, got %v", out["Name"])
	}
	if out["value"] != float64(42) {
		t.Errorf("expected 42, got %v", out["Value"])
	}
	if _, ok := out["hidden"]; ok {
		t.Error("unexported field 'hidden' should not appear")
	}
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
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["__class__"] != "json.jsonStruct" {
		t.Errorf("expected __class__ json.jsonStruct, got %v", out["__class__"])
	}
	if out["name"] != "test" {
		t.Errorf("expected test, got %v", out["name"])
	}
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
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if m["name"] != "a" {
		t.Errorf("expected 'a', got %v", m["Name"])
	}
	// The cycle reference should produce a string marker, not a recursive struct
	next, ok := m["next"].(map[string]any)
	if !ok {
		t.Fatalf("expected map for Next, got %T", m["Next"])
	}
	if next["name"] != "b" {
		t.Errorf("expected 'b', got %v", next["Name"])
	}
	// b.Next should be a cycle marker string
	cycle, ok := next["next"].(string)
	if !ok {
		t.Fatalf("expected string cycle marker, got %T", next["Next"])
	}
	if len(cycle) < 3 {
		t.Errorf("expected meaningful cycle marker, got %q", cycle)
	}
}

func TestAsJSONs(t *testing.T) {
	v := map[string]any{"x": 1}
	s := AsJSONStr(v)
	want := "{\n  \"x\": 1\n}"
	if s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}

func TestAsJSONChanFunc(t *testing.T) {
	ch := make(chan int)
	fn := func() {}
	v := []any{ch, fn}
	b, err := json.Marshal(AsJSON(v))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `[null,null]` {
		t.Errorf("expected [null,null], got %s", string(b))
	}
}
