// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/asjson"
)

func TestNodeNilSafety(t *testing.T) {
	var n *Node
	assert.Zero(t, n.Parent())
	n.setParent(&Node{})
	assert.Equal(t, "", n.Text())
	assert.Equal(t, 0, n.Line())
	assert.Zero(t, n.Path())
	b, err := json.Marshal(asjson.AsJSON(n))
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestNodeParent(t *testing.T) {
	parent := &Node{}
	child := &Node{}
	child.setParent(parent)
	assert.Equal(t, parent, child.Parent())
	assert.Zero(t, parent.Parent())
}

func TestNodeText(t *testing.T) {
	n := &Node{
		ParseInfo: &ParseInfo{
			Source: "hello world",
			Start:  6,
			Mark:   11,
		},
	}
	assert.Equal(t, "world", n.Text())
}

func TestNodeTextNilPos(t *testing.T) {
	n := &Node{}
	assert.Equal(t, "", n.Text())
}

func TestNodeLine(t *testing.T) {
	n := &Node{
		ParseInfo: &ParseInfo{Line: 3},
	}
	assert.Equal(t, 3, n.Line())
}

func TestNodeLineNilPos(t *testing.T) {
	n := &Node{}
	assert.Equal(t, 0, n.Line())
}

func TestNodePath(t *testing.T) {
	root := &Node{Ast: "root"}
	mid := &Node{Ast: "mid"}
	leaf := &Node{Ast: "leaf"}
	mid.setParent(root)
	leaf.setParent(mid)
	path := leaf.Path()
	assert.Equal(t, 3, len(path))
	assert.Equal(t, root, path[0])
	assert.Equal(t, mid, path[1])
	assert.Equal(t, leaf, path[2])
}

func TestNodePathSingle(t *testing.T) {
	n := &Node{}
	path := n.Path()
	assert.Equal(t, 1, len(path))
	assert.Equal(t, n, path[0])
}

func TestNodePathNil(t *testing.T) {
	var n *Node
	assert.Zero(t, n.Path())
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
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	assert.Equal(t, "peg.Node", out["__class__"])
	// FIXME
	//if out["ast"] != nil {
	//	t.Errorf("expected no ast, got %v", out["ast"])
	//}
	parseInfo, ok := out["parse_info"].(map[string]any)
	assert.True(t, ok, "expected pos map")
	assert.Equal(t, "src", parseInfo["source"])
}

func TestNodeMarshalJSONNil(t *testing.T) {
	var n *Node
	b, err := json.Marshal(asjson.AsJSON(n))
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestNodeMarshalJSONAstMap(t *testing.T) {
	n := &Node{
		Ast: map[string]any{
			"key": "value",
			"num": float64(42),
		},
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	ast, ok := out["ast"].(map[string]any)
	assert.True(t, ok, "expected ast to be object")
	assert.Equal(t, "value", ast["key"])
	assert.Equal[any](t, float64(42), ast["num"])
}

func TestNodeMarshalJSONAstSlice(t *testing.T) {
	n := &Node{
		Ast: []any{"a", "b", 3},
	}
	b, err := json.Marshal(asjson.AsJSON(n))
	assert.NoError(t, err)
	var out map[string]any
	err = json.Unmarshal(b, &out)
	assert.NoError(t, err)
	ast, ok := out["ast"].([]any)
	assert.True(t, ok, "expected ast to be array")
	assert.Equal(t, 3, len(ast))
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
	assert.Equal(t, "test.ebnf", pi.Source)
	assert.Equal(t, "number", pi.Rule)
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
	assert.Equal(t, 1, len(children))
	assert.Equal(t, child, children[0])
}
