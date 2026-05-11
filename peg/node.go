package peg

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	parent   *Node
	Ast      any
	Pos      *ParseInfo
	Children []*Node
}

func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	return n.parent
}

func (n *Node) SetParent(p *Node) {
	if n != nil {
		n.parent = p
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
	for cur := n; cur != nil; cur = cur.parent {
		ancestors = append(ancestors, cur)
	}
	for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
	}
	return ancestors
}

func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	cp := *n
	cp.parent = nil
	cp.Children = append([]*Node(nil), n.Children...)
	return &cp
}

// marshalValue recursively converts a Go value to a JSON-safe form.
func marshalValue(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		return v
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		out := make([]any, 0, n)
		for i := range n {
			out = append(out, marshalValue(rv.Index(i).Interface()))
		}
		return out
	case reflect.Map:
		out := make(map[string]any, rv.Len())
		for _, key := range rv.MapKeys() {
			k := fmt.Sprint(key.Interface())
			out[k] = marshalValue(rv.MapIndex(key).Interface())
		}
		return out
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return marshalValue(rv.Elem().Interface())
	case reflect.Struct:
		if m, ok := v.(json.Marshaler); ok {
			raw, _ := m.MarshalJSON()
			var out any
			_ = json.Unmarshal(raw, &out)
			return out
		}
		out := make(map[string]any)
		t := rv.Type()
		for i := range rv.NumField() {
			name := t.Field(i).Name
			if !t.Field(i).IsExported() {
				continue
			}
			out[name] = marshalValue(rv.Field(i).Interface())
		}
		return out
	default:
		return fmt.Sprint(v)
	}
}

func (n *Node) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	out := make(map[string]any)
	out["__class__"] = "Node"
	if n.Ast != nil {
		out["ast"] = marshalValue(n.Ast)
	}
	if n.Pos != nil {
		out["pos"] = marshalValue(n.Pos)
	}
	if len(n.Children) > 0 {
		children := make([]any, 0, len(n.Children))
		for _, child := range n.Children {
			if child == nil {
				children = append(children, nil)
			} else {
				children = append(children, child)
			}
		}
		out["children"] = children
	}
	return json.Marshal(out)
}
