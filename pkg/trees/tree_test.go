// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/asjson"
)

func text(s string) string      { return s }
func seq(items ...any) *treeSeq { return &treeSeq{Items: items} }

func TestFoldBottom(t *testing.T) {
	result := Fold(&typeBottomTree{})
	_, ok := result.(*typeBottomTree)
	assert.True(t, ok, "expected Bottom, got %T", result)
}

func TestFoldGoNil(t *testing.T) {
	result := Fold(nil)
	assert.True(t, result == nil, "expected nil, got %T", result)
}

func TestFoldText(t *testing.T) {
	result := Fold("hello")
	s, ok := result.(string)
	assert.True(t, ok, "expected string, got %T", result)
	assert.Equal(t, "hello", s)
}

func TestFoldBool(t *testing.T) {
	result := Fold(true)
	b, ok := result.(bool)
	assert.True(t, ok, "expected bool, got %T", result)
	assert.Equal(t, true, b)
}

func TestFoldNumber(t *testing.T) {
	result := Fold(42.5)
	f, ok := result.(float64)
	assert.True(t, ok, "expected float64, got %T", result)
	assert.Equal(t, 42.5, f)
}

func TestFoldSeqToSeq(t *testing.T) {
	result := Fold(seq(text("a"), text("b"), text("c")))
	l, ok := result.([]any)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 3, len(l))
	assert.Equal(t, "a", l[0].(string), "expected 'a'")
}

func TestFoldListToList(t *testing.T) {
	result := Fold([]any{text("a"), text("b"), text("c")})
	l, ok := result.([]any)
	assert.True(t, ok, "expected list, got %T", result)
	assert.Equal(t, 3, len(l))
	assert.Equal(t, "a", l[0].(string), "expected 'a'")
}

func TestFoldNamedToMap(t *testing.T) {
	result := Fold(NamedTree("x", text("hello")))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.NotZero(t, m["x"], "expected key 'x'")
	assert.Equal(t, "hello", m["x"].(string), "expected 'hello'")
}

func TestFoldOverride(t *testing.T) {
	result := Fold(OverrideTree(text("result")))
	s, ok := result.(string)
	assert.True(t, ok, "expected Text, got %T", result)
	assert.Equal(t, "result", s, "expected 'result'")
}

func TestFoldMultipleNamed(t *testing.T) {
	result := Fold(seq(
		NamedTree("a", text("1")),
		NamedTree("b", text("2")),
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, "1", m["a"].(string), "expected '1'")
	assert.Equal(t, "2", m["b"].(string), "expected '2'")
}

func TestFoldNamedAccumulates(t *testing.T) {
	result := Fold(seq(
		NamedTree("x", text("a")),
		NamedTree("x", text("b")),
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, "a", m["x"].([]any)[0].(string), "expected 'a'")
	assert.Equal(t, "b", m["x"].([]any)[1].(string), "expected 'b'")
}

func TestFoldNamedAsList(t *testing.T) {
	result := Fold(NamedTreeSeq("items", text("x")))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, 1, len(m["items"].([]any)))
	assert.Equal(t, "x", m["items"].([]any)[0].(string), "expected 'x'")
}

func TestFoldNamedAsListAccumulates(t *testing.T) {
	result := Fold(seq(
		NamedTreeSeq("items", text("a")),
		NamedTreeSeq("items", text("b")),
	))
	m, ok := result.(map[string]any)
	assert.True(t, ok, "expected MapNode, got %T", result)
	assert.Equal(t, 2, len(m["items"].([]any)))
	assert.Equal(t, "a", m["items"].([]any)[0].(string), "expected 'a'")
	assert.Equal(t, "b", m["items"].([]any)[1].(string), "expected 'b'")
}

func TestFoldOverrideWins(t *testing.T) {
	result := Fold(seq(
		NamedTree("x", text("ignored")),
		text("also ignored"),
		OverrideTree(text("result")),
	))
	s, ok := result.(string)
	assert.True(t, ok, "expected string, got %v", result)
	assert.Equal(t, "result", s, "expected 'result'")
}

func TestFoldOverrideAsList(t *testing.T) {
	result := Fold(seq(
		OverrideTreeSeq(text("a")),
		OverrideTreeSeq(text("b")),
	))
	l, ok := result.([]any)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 2, len(l))
}

func TestFoldNestedNamed(t *testing.T) {
	result := Fold(NamedTree("x", seq(
		NamedTree("a", text("1")),
		NamedTree("b", text("2")),
	)))
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
	assert.Equal(t, "hello", treeToJSON("hello"))
}

func TestNumberAsJSON(t *testing.T) {
	result := treeToJSON(42.5)
	assert.Equal(t, 42.5, result)
}

func TestNodeAsJSONTree(t *testing.T) {
	n := &Node{TypeName: "expr", Tree: text("42")}
	result := asjson.AsJSONStr(n)
	want := "{\n  \"__class__\": \"expr\",\n  \"ast\": \"42\"\n}"
	assert.Equal(t, want, result)
}

func TestFoldRuleNode(t *testing.T) {
	result := Fold(&Node{TypeName: "expr", Tree: text("42")})
	r, ok := result.(*Node)
	assert.True(t, ok, "expected RuleNode, got %T", result)
	assert.Equal(t, "expr", r.TypeName, "expected 'expr'")
}
