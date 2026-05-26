// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
	"weak"

	"github.com/neogeny/ogopego/util"
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
	return fmt.Sprintf("%+v", n)
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
	children := make([]*Node, 0)

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
			for _, item := range val.Entries() {
				dfs(item)
			}
			return
		default:
			return
		}
	}
	dfs(n.Ast)
	pub := util.PubMapOf(n)
	if m, ok := pub.(OrderedMap); ok {
		for _, obj := range m.Entries() {
			dfs(obj)
		}
	}
	return children
}
