package test

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/tool"
	"github.com/neogeny/ogopego/pkg/trees"
)

func TestModelReprOutput(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	code := tool.ModelRepr(*g, "calc")
	assert.NotZero(t, code, "empty output")
}

func TestPairFromTree(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	assert.NoError(t, err, "parse")

	n, ok := tree.(*trees.Node)
	assert.True(t, ok, "expected *trees.Node, got %T", tree)
	assert.Equal(t, "Pair", n.TypeName)
	m, ok := n.Tree.(*trees.MapNode)
	assert.True(t, ok, "expected MapNode, got %T", n.Tree)
	assert.Equal(t, 2, len(m.Entries), "expected 2 entries")
	key := m.Entries["key"].(*trees.Text).Value
	val := m.Entries["val"].(*trees.Text).Value
	assert.Equal(t, "abc", key, "key")
	assert.Equal(t, "123", val, "val")
}

func TestModelReprTypedRef(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
		NUMBER::Num = /\d+/
	`), nil)
	code := tool.ModelRepr(*g, "calc")
	assert.True(t, strings.Contains(code, "Child *Pair"), "expected Start.Child *Pair field")
	assert.True(t, strings.Contains(code, "PairFromTree("), "expected PairFromTree call in StartFromTree")
	assert.True(t, strings.Contains(code, "Value any"), "expected Num.Value any field")
}

func TestTypedRefFromTree(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`), nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	assert.NoError(t, err, "parse")
	// start::Start wraps: Node{Start, MapNode{child: Node{Pair, MapNode{key, val}}}}
	ns, ok := tree.(*trees.Node)
	assert.True(t, ok, "expected *trees.Node, got %T", tree)
	assert.Equal(t, "Start", ns.TypeName)
	m, ok := ns.Tree.(*trees.MapNode)
	assert.True(t, ok, "expected MapNode, got %T", ns.Tree)
	childTree := m.Entries["child"]
	np, ok := childTree.(*trees.Node)
	assert.True(t, ok, "expected *trees.Node for child, got %T", childTree)
	assert.Equal(t, "Pair", np.TypeName)
	pm, ok := np.Tree.(*trees.MapNode)
	assert.True(t, ok, "expected MapNode for Pair, got %T", np.Tree)
	key := pm.Entries["key"].(*trees.Text).Value
	val := pm.Entries["val"].(*trees.Text).Value
	assert.Equal(t, "abc", key, "key")
	assert.Equal(t, "123", val, "val")
}
