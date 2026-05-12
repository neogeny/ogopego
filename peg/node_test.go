package peg

import (
	"encoding/json"
	"testing"
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
	if c := n.Clone(); c != nil {
		t.Error("expected nil clone")
	}
	b, err := n.MarshalJSON()
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
		Pos: &ParseInfo{
			Source: "hello world",
			Pos:    6,
			EndPos: 11,
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
		Pos: &ParseInfo{Line: 3},
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

func TestNodeClone(t *testing.T) {
	parent := &Node{}
	n := &Node{
		Ast: "value",
		Pos: &ParseInfo{Line: 5, Source: "src"},
		Children: []*Node{
			{Ast: "child"},
		},
	}
	n.setParent(parent)
	c := n.Clone()
	if c.Parent() != nil {
		t.Error("clone should have nil parent")
	}
	if c.Ast != n.Ast {
		t.Error("clone should share ast")
	}
	if c.Pos != n.Pos {
		t.Error("clone should share pos pointer")
	}
	if len(c.Children) != 1 {
		t.Fatal("expected 1 child")
	}
	if c.Children[0] != n.Children[0] {
		t.Error("clone should share child pointers")
	}
	c.Children[0] = &Node{Ast: "replaced"}
	if n.Children[0].Ast != "child" {
		t.Error("clone children slice should be independent")
	}
}

func TestNodeCloneNil(t *testing.T) {
	var n *Node
	if c := n.Clone(); c != nil {
		t.Error("expected nil clone")
	}
}

func TestNodeMarshalJSON(t *testing.T) {
	n := &Node{
		Ast: "simple",
		Pos: &ParseInfo{Source: "src", Rule: "rule", Pos: 0, EndPos: 6, Line: 1},
	}
	b, err := n.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["__class__"] != "Node" {
		t.Errorf("expected __class__ Node, got %v", out["__class__"])
	}
	if out["Ast"] != "simple" {
		t.Errorf("expected Ast 'simple', got %v", out["Ast"])
	}
	pos, ok := out["Pos"].(map[string]any)
	if !ok {
		t.Fatal("expected Pos map")
	}
	if pos["Source"] != "src" {
		t.Errorf("expected Source 'src', got %v", pos["Source"])
	}
	if out["Children"] != nil {
		t.Error("expected nil Children")
	}
}

func TestNodeMarshalJSONNil(t *testing.T) {
	var n *Node
	b, err := n.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Errorf("expected null, got %s", b)
	}
}

func TestNodeMarshalJSONWithChildren(t *testing.T) {
	child := &Node{Ast: "child"}
	parent := &Node{
		Ast:      "parent",
		Children: []*Node{child},
	}
	child.setParent(parent)
	b, err := parent.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	children, ok := out["Children"].([]any)
	if !ok {
		t.Fatal("expected Children array")
	}
	if len(children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(children))
	}
	c, ok := children[0].(map[string]any)
	if !ok {
		t.Fatal("expected child to be object")
	}
	if c["Ast"] != "child" {
		t.Errorf("expected child Ast 'child', got %v", c["Ast"])
	}
}

func TestNodeMarshalJSONAstMap(t *testing.T) {
	n := &Node{
		Ast: map[string]any{
			"key": "value",
			"num": float64(42),
		},
	}
	b, err := n.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	ast, ok := out["Ast"].(map[string]any)
	if !ok {
		t.Fatal("expected Ast to be object")
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
	b, err := n.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	ast, ok := out["Ast"].([]any)
	if !ok {
		t.Fatal("expected Ast to be array")
	}
	if len(ast) != 3 {
		t.Fatalf("expected 3 items, got %d", len(ast))
	}
}

func TestParseInfo(t *testing.T) {
	pi := &ParseInfo{
		Source:  "test.ebnf",
		Rule:    "number",
		Pos:     10,
		EndPos:  13,
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
