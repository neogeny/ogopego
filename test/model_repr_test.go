package test

import (
	"strings"
	"testing"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/tool"
	"github.com/neogeny/ogopego/trees"
)

func TestModelReprOutput(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	code := tool.ModelRepr(*g, "calc")
	if code == "" {
		t.Fatal("empty output")
	}
}

func TestPairFromTree(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	n, ok := tree.(*trees.Node)
	if !ok {
		t.Fatalf("expected *trees.Node, got %T", tree)
	}
	if n.TypeName != "Pair" {
		t.Fatalf("expected TypeName Pair, got %q", n.TypeName)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", n.Tree)
	}
	if len(m.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m.Entries))
	}
	key := m.Entries["key"].(*trees.Text).Value
	val := m.Entries["val"].(*trees.Text).Value
	if key != "abc" || val != "123" {
		t.Fatalf("expected (abc,123), got (%q,%q)", key, val)
	}
}

func TestModelReprTypedRef(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
		NUMBER::Num = /\d+/
	`), nil)
	code := tool.ModelRepr(*g, "calc")
	if !strings.Contains(code, "Child *Pair") {
		t.Error("expected Start.Child *Pair field")
	}
	if !strings.Contains(code, "PairFromTree(") {
		t.Error("expected PairFromTree call in StartFromTree")
	}
	if !strings.Contains(code, "Value any") {
		t.Error("expected Num.Value any field")
	}
}

func TestTypedRefFromTree(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`), nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// start::Start wraps: Node{Start, MapNode{child: Node{Pair, MapNode{key, val}}}}
	ns, ok := tree.(*trees.Node)
	if !ok {
		t.Fatalf("expected *trees.Node, got %T", tree)
	}
	if ns.TypeName != "Start" {
		t.Fatalf("expected TypeName Start, got %q", ns.TypeName)
	}
	m, ok := ns.Tree.(*trees.MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", ns.Tree)
	}
	childTree := m.Entries["child"]
	np, ok := childTree.(*trees.Node)
	if !ok {
		t.Fatalf("expected *trees.Node for child, got %T", childTree)
	}
	if np.TypeName != "Pair" {
		t.Fatalf("expected TypeName Pair, got %q", np.TypeName)
	}
	pm, ok := np.Tree.(*trees.MapNode)
	if !ok {
		t.Fatalf("expected MapNode for Pair, got %T", np.Tree)
	}
	key := pm.Entries["key"].(*trees.Text).Value
	val := pm.Entries["val"].(*trees.Text).Value
	if key != "abc" || val != "123" {
		t.Fatalf("expected (abc, 123), got (%q, %q)", key, val)
	}
}
