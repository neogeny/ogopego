package peg

import (
	"encoding/json"
	"reflect"
	"weak"

	"github.com/davecgh/go-spew/spew"
)

// ParseInfo holds parse position information.
type ParseInfo struct {
	Source  string
	Rule    string
	Pos     int
	EndPos  int
	Line    int
	EndLine int
}

// Node is an AST node produced by parsing.
type Node struct {
	parent   weak.Pointer[Node]
	Ast      any
	Pos      *ParseInfo
	Children []*Node
}

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

func (n *Node) Text() string {
	if n == nil || n.Pos == nil {
		return ""
	}
	return n.Pos.Source[n.Pos.Pos:n.Pos.EndPos]
}

func (n *Node) Line() int {
	if n == nil || n.Pos == nil {
		return 0
	}
	return n.Pos.Line
}

func (n *Node) AsStr() string {
	return spew.Sdump(n)
}

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

func (n *Node) __pub__() map[string]any {
	if n == nil {
		return nil
	}
	out := make(map[string]any)
	v := reflect.ValueOf(n).Elem()
	t := v.Type()
	for i := range v.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		out[f.Name] = v.Field(i).Interface()
	}
	return out
}

func (n *Node) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	pub := n.__pub__()
	pub["__class__"] = "Node"
	return json.Marshal(pub)
}
