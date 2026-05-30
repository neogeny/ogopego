// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"encoding/json"
	"testing"

	"github.com/neogeny/ogopego/pkg/asjson"
)

func TestNodeNilSafety(t *testing.T) {
	var n *Node
	if p := n.Parent(); p != nil {
		t.Error("expected nil parent")
	}
	n.setParent(&Node{})
	if s := n.Text(); s != "" {
		t.Errorf("expected empty text, got %q", s)
	}
	if l := n.Line(); l != 0 {
		t.Errorf("expected 0 line, got %d", l)
	}
	if p := n.Path(); p != nil {
		t.Error("expected nil path")
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Errorf("expected null, got %s", b)
	}
}

func TestNodeParent(t *testing.T) {
	parent := &Node{}
	child := &Node{}
	child.setParent(parent)
	if child.Parent() != parent {
		t.Error("expected parent to match")
	}
	if parent.Parent() != nil {
		t.Error("expected root parent to be nil")
	}
}

func TestNodeText(t *testing.T) {
	n := &Node{
		ParseInfo: &ParseInfo{
			Source: "hello world",
			Start:  6,
			Mark:   11,
		},
	}
	if s := n.Text(); s != "world" {
		t.Errorf("expected 'world', got %q", s)
	}
}

func TestNodeTextNilPos(t *testing.T) {
	n := &Node{}
	if s := n.Text(); s != "" {
		t.Errorf("expected empty text, got %q", s)
	}
}

func TestNodeLine(t *testing.T) {
	n := &Node{
		ParseInfo: &ParseInfo{Line: 3},
	}
	if l := n.Line(); l != 3 {
		t.Errorf("expected line 3, got %d", l)
	}
}

func TestNodeLineNilPos(t *testing.T) {
	n := &Node{}
	if l := n.Line(); l != 0 {
		t.Errorf("expected 0, got %d", l)
	}
}

func TestNodePath(t *testing.T) {
	root := &Node{Ast: "root"}
	mid := &Node{Ast: "mid"}
	leaf := &Node{Ast: "leaf"}
	mid.setParent(root)
	leaf.setParent(mid)
	path := leaf.Path()
	if len(path) != 3 {
		t.Fatalf("expected 3 ancestors, got %d", len(path))
	}
	if path[0] != root {
		t.Error("expected root first")
	}
	if path[1] != mid {
		t.Error("expected mid second")
	}
	if path[2] != leaf {
		t.Error("expected leaf third")
	}
}

func TestNodePathSingle(t *testing.T) {
	n := &Node{}
	path := n.Path()
	if len(path) != 1 || path[0] != n {
		t.Errorf("expected [self], got %v", path)
	}
}

func TestNodePathNil(t *testing.T) {
	var n *Node
	if p := n.Path(); p != nil {
		t.Error("expected nil path")
	}
}

func TestNodeMarshalJSON(t *testing.T) {
	n := &Node{
		Ast: "simple",
		ParseInfo: &ParseInfo{
			Source: "src",
			Rule:   "rule",
			Start:  0,
			Mark:   6,
			Line:   1,
		},
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["__class__"] != "peg.Node" {
		t.Errorf("expected __class__ peg.Node, got %v", out["__class__"])
	}
	// FIXME
	//if out["ast"] != nil {
	//	t.Errorf("expected no ast, got %v", out["ast"])
	//}
	parseInfo, ok := out["parse_info"].(map[string]any)
	if !ok {
		t.Errorf("expected pos map")
	}
	if parseInfo["source"] != "src" {
		t.Errorf("expected Source 'src', got %v", parseInfo["Source"])
	}
}

func TestNodeMarshalJSONNil(t *testing.T) {
	var n *Node
	b, err := json.Marshal(asjson.AsJSON(n))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Errorf("expected null, got %s", b)
	}
}

func TestNodeMarshalJSONAstMap(t *testing.T) {
	n := &Node{
		Ast: map[string]any{
			"key": "value",
			"num": float64(42),
		},
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	ast, ok := out["ast"].(map[string]any)
	if !ok {
		t.Fatal("expected ast to be object")
	}
	if ast["key"] != "value" {
		t.Errorf("expected value, got %v", ast["key"])
	}
	if ast["num"] != float64(42) {
		t.Errorf("expected 42, got %v", ast["num"])
	}
}

func TestNodeMarshalJSONAstSlice(t *testing.T) {
	n := &Node{
		Ast: []any{"a", "b", 3},
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	ast, ok := out["ast"].([]any)
	if !ok {
		t.Fatal("expected ast to be array")
	}
	if len(ast) != 3 {
		t.Fatalf("expected 3 items, got %d", len(ast))
	}
}

func TestParseInfo(t *testing.T) {
	pi := &ParseInfo{
		Source:  "test.ebnf",
		Rule:    "number",
		Start:   10,
		Mark:    13,
		Line:    2,
		EndLine: 2,
	}
	if pi.Source != "test.ebnf" {
		t.Errorf("expected 'test.ebnf', got %q", pi.Source)
	}
	if pi.Rule != "number" {
		t.Errorf("expected 'number', got %q", pi.Rule)
	}
}

func TestNodeChildren(t *testing.T) {
	child := &Node{Ast: "child"}
	parent := &Node{
		Ast: map[string]any{
			"key":   "value",
			"child": child,
		},
	}
	children := parent.Children()
	if len(children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(children))
	}
	if children[0] != child {
		t.Error("expected child")
	}
}
