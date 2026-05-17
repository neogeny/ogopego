// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"encoding/json"
	"weak"

	"github.com/davecgh/go-spew/spew"
	asjson "github.com/neogeny/ogopego/json"
)

// ParseInfo holds parse position information.
type ParseInfo struct {
	Source  string
	Rule    string
	Start   int
	Mark    int
	Line    int
	EndLine int
}

// Node is an AST node produced by parsing.
type Node struct {
	asjson.AsJSONBase
	parent    weak.Pointer[Node]
	Ast       any
	ParseInfo *ParseInfo
	children  []*Node
}

// Parent returns the parent node of the current node.
func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	return n.parent.Value()
}

func (n *Node) setParent(p *Node) {
	if n != nil {
		n.parent = weak.Make(p)
	}
}

// Text returns the text segment matched by this node.
func (n *Node) Text() string {
	if n == nil || n.ParseInfo == nil {
		return ""
	}
	return n.ParseInfo.Source[n.ParseInfo.Start:n.ParseInfo.Mark]
}

// Line returns the starting line number of this node.
func (n *Node) Line() int {
	if n == nil || n.ParseInfo == nil {
		return 0
	}
	return n.ParseInfo.Line
}

// AsStr returns a string representation of the node for debugging.
func (n *Node) AsStr() string {
	return spew.Sdump(n)
}

// Path returns the path of ancestors from the root to this node.
func (n *Node) Path() []*Node {
	if n == nil {
		return nil
	}
	var ancestors []*Node
	for cur := n; cur != nil; cur = cur.Parent() {
		ancestors = append(ancestors, cur)
	}
	for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
	}
	return ancestors
}

// PubMap returns an ordered map of the node's public fields.
func (n *Node) PubMap() *asjson.OrderedMap {
	if n == nil {
		return nil
	}
	pub := n.AsJSONBase.PubMapOf(n)
	if val, ok := pub.Get("parse_info"); ok {
		if val == nil {
			pub.Delete("parse_info")
		}
	}
	if len(pub.Keys()) > 1 {
		pub.Delete("Ast")
	}
	return pub
}

// AsJSON returns a JSON-compatible representation of the node.
func (n *Node) AsJSON() any {
	if n == nil {
		return nil
	}
	return n.AsJSONBase.AsJSONOf(n)
}

// AsJSONStr returns a JSON string representation of the node.
func (n *Node) AsJSONStr() string { return n.AsJSONStrOf(n) }

// MarshalJSON marshals the node to JSON.
func (n *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.AsJSON())
}

// Children returns the children of the current node.
func (n *Node) Children() []*Node {
	if n == nil {
		return nil
	}
	if n.children == nil {
		n.children = n.getChildren()
	}
	return n.children
}

func (n *Node) getChildren() []*Node {
	if n == nil {
		return nil
	}
	pub := n.PubMap()
	children := make([]*Node, 0, len(pub.Keys()))

	dfs := func(obj any) {}
	dfs = func(obj any) {
		if obj == nil {
			return
		}
		switch val := obj.(type) {
		case *Node:
			val.setParent(n)
			children = append(children, val)
			return
		case []any:
			for _, item := range val {
				dfs(item)
			}
			return
		case map[string]any:
			for _, item := range val {
				dfs(item)
			}
			return
		case OrderedMap:
			for _, item := range val.Values() {
				dfs(item)
			}
			return
		default:
			return
		}
	}
	dfs(n.Ast)
	for obj := range pub.Values() {
		dfs(obj)
	}
	return children
}
