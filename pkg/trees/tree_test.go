// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/asjson"
)

func text(s string) *Text   { return &Text{Value: s} }
func seq(items ...any) *Seq { return &Seq{Items: items} }

func TestFoldBottom(t *testing.T) {
	result := Fold(&Bottom{})
	_, ok := result.(*Bottom)
	assert.True(t, ok, "expected Bottom, got %T", result)
}

func TestFoldGoNil(t *testing.T) {
	result := Fold(nil)
	assert.True(t, result == nil, "expected nil, got %T", result)
}

func TestFoldText(t *testing.T) {
	result := Fold(text("hello"))
	txt, ok := result.(*Text)
	assert.True(t, ok, "expected Text, got %T", result)
	assert.Equal(t, "hello", txt.Value)
}

func TestFoldBool(t *testing.T) {
	result := Fold(true)
	b, ok := result.(bool)
	assert.True(t, ok, "expected bool, got %T", result)
	assert.Equal(t, true, b)
}

func TestFoldNumber(t *testing.T) {
	result := Fold(&Number{Value: 42.5})
	n, ok := result.(*Number)
	assert.True(t, ok, "expected Number, got %T", result)
	assert.Equal(t, 42.5, n.Value)
}

func TestFoldSeqToSeq(t *testing.T) {
	result := Fold(seq(text("a"), text("b"), text("c")))
	l, ok := result.([]any)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 3, len(l))
	assert.Equal(t, "a", l[0].(*Text).Value, "expected 'a'")
}

func TestFoldListToList(t *testing.T) {
	result := Fold([]any{text("a"), text("b"), text("c")})
	l, ok := result.([]any)
	assert.True(t, ok, "expected list, got %T", result)
	assert.Equal(t, 3, len(l))
	assert.Equal(t, "a", l[0].(*Text).Value, "expected 'a'")
}

func TestFoldNamedToMap(t *testing.T) {
	result := Fold(&Named{Name: "x", Value: text("hello")})
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.NotZero(t, m["x"], "expected key 'x'")
	assert.Equal(t, "hello", m["x"].(*Text).Value, "expected 'hello'")
}

func TestFoldOverride(t *testing.T) {
	result := Fold(&Override{Value: text("result")})
	txt, ok := result.(*Text)
	assert.True(t, ok, "expected Text, got %T", result)
	assert.Equal(t, "result", txt.Value, "expected 'result'")
}

func TestFoldMultipleNamed(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "a", Value: text("1")},
		&Named{Name: "b", Value: text("2")},
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, "1", m["a"].(*Text).Value, "expected '1'")
	assert.Equal(t, "2", m["b"].(*Text).Value, "expected '2'")
}

func TestFoldNamedAccumulates(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "x", Value: text("a")},
		&Named{Name: "x", Value: text("b")},
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, "a", m["x"].([]any)[0].(*Text).Value, "expected 'a'")
	assert.Equal(t, "b", m["x"].([]any)[1].(*Text).Value, "expected 'b'")
}

func TestFoldNamedAsList(t *testing.T) {
	result := Fold(&NamedAsList{Name: "items", Value: text("x")})
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, 1, len(m["items"].([]any)))
	assert.Equal(t, "x", m["items"].([]any)[0].(*Text).Value, "expected 'x'")
}

func TestFoldNamedAsListAccumulates(t *testing.T) {
	result := Fold(seq(
		&NamedAsList{Name: "items", Value: text("a")},
		&NamedAsList{Name: "items", Value: text("b")},
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, 2, len(m["items"].([]any)))
	assert.Equal(t, "a", m["items"].([]any)[0].(*Text).Value, "expected 'a'")
	assert.Equal(t, "b", m["items"].([]any)[1].(*Text).Value, "expected 'b'")
}

func TestFoldOverrideWins(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "x", Value: text("ignored")},
		text("also ignored"),
		&Override{Value: text("result")},
	))
	txt, ok := result.(*Text)
	assert.True(t, ok, "expected Text, got %T", result)
	assert.Equal(t, "result", txt.Value, "expected 'result'")
}

func TestFoldOverrideAsList(t *testing.T) {
	result := Fold(seq(
		&OverrideAsList{Value: text("a")},
		&OverrideAsList{Value: text("b")},
	))
	l, ok := result.([]any)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 2, len(l))
}

func TestFoldNestedNamed(t *testing.T) {
	result := Fold(&Named{
		Name: "x",
		Value: seq(
			&Named{Name: "a", Value: text("1")},
			&Named{Name: "b", Value: text("2")},
		),
	})
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	_, exists := m["x"]
	assert.True(t, exists, "expected key 'x'")
	_, exists = m["a"]
	assert.True(t, exists, "expected key 'a'")
	_, exists = m["b"]
	assert.True(t, exists, "expected key 'b'")
}

func TestFoldSeqWithNil(t *testing.T) {
	result := Fold(seq(text("a"), nil, text("b")))
	l, ok := result.([]any)
	assert.True(t, ok, "expected list, got %T", result)
	assert.Equal(t, 2, len(l))
}

func TestTextAsJSON(t *testing.T) {
	var aj asjson.AsJSONMixin = &Text{Value: "hello"}
	result := asjson.AsJSONStr(aj.As_JSON_())
	assert.Equal(t, `"hello"`, result)
}

func TestNumberAsJSON(t *testing.T) {
	var aj asjson.AsJSONMixin = &Number{Value: 42.5}
	result := asjson.AsJSONStr(aj.As_JSON_())
	assert.Equal(t, "42.5", result)
}

func TestNodeAsJSONTree(t *testing.T) {
	n := &Node{TypeName: "expr", Tree: text("42")}
	var aj asjson.AsJSONMixin = n
	result := asjson.AsJSONStr(aj.As_JSON_())
	want := "{\n  \"__class__\": \"expr\",\n  \"ast\": \"42\"\n}"
	assert.Equal(t, want, result)
}

func TestFoldRuleNode(t *testing.T) {
	result := Fold(&Node{TypeName: "expr", Tree: text("42")})
	r, ok := result.(*Node)
	assert.True(t, ok, "expected RuleNode, got %T", result)
	assert.Equal(t, "expr", r.TypeName, "expected 'expr'")
}
