// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"fmt"
	"maps"

	"github.com/neogeny/ogopego/pkg/asjson"
)

var _ asjson.AsJSONMixin = (*Node)(nil)

// / treeToJSONStr
func treeToJSON(t any) any {
	if t == nil {
		return nil
	}
	if _, ok := t.(asjson.AsJSONMixin); !ok {
		panic("TreeToJSON: not an AsJSONMixin: " + fmt.Sprintf("%T", t))
	}
	switch v := t.(type) {
	case *Text:
		return v.Value
	case *Number:
		return v.Value
	case *Bool:
		return v.Value
	case *Nil:
		return nil
	case *Bottom:
		return nil
	case *TrueValue:
		return true
	case *FalseValue:
		return false
	case *NullValue:
		return nil
	case *Seq:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = treeToJSON(item)
		}
		return items
	case *Array:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = treeToJSON(item)
		}
		return items
	case *MapNode:
		out := make(map[string]any, len(v.Entries))
		for k, val := range v.Entries {
			out[k] = treeToJSON(val)
		}
		return out
	case *Named:
		return map[string]any{v.Name: treeToJSON(v.Value)}
	case *NamedAsList:
		return map[string]any{v.Name: treeToJSON(v.Value)}
	case *Override:
		return treeToJSON(v.Value)
	case *OverrideAsList:
		return treeToJSON(v.Value)
	case *Node:
		child := treeToJSON(v.Tree)
		if m, ok := child.(map[string]any); ok {
			if _, has := m["__class__"]; !has {
				out := make(map[string]any, len(m)+1)
				out["__class__"] = v.TypeName
				maps.Copy(out, m)
				return out
			}
		}
		return map[string]any{"__class__": v.TypeName, "ast": child}
	default:
		return t
	}
}

// As_JSON_ implementations for each concrete tree type.

func (t *Text) As_JSON_() any           { return treeToJSON(t) }
func (n *Number) As_JSON_() any         { return treeToJSON(n) }
func (b *Bool) As_JSON_() any           { return treeToJSON(b) }
func (*Nil) As_JSON_() any              { return nil }
func (*Bottom) As_JSON_() any           { return nil }
func (*TrueValue) As_JSON_() any        { return true }
func (*FalseValue) As_JSON_() any       { return false }
func (*NullValue) As_JSON_() any        { return nil }
func (s *Seq) As_JSON_() any            { return treeToJSON(s) }
func (a *Array) As_JSON_() any          { return treeToJSON(a) }
func (m *MapNode) As_JSON_() any        { return treeToJSON(m) }
func (n *Named) As_JSON_() any          { return treeToJSON(n) }
func (n *NamedAsList) As_JSON_() any    { return treeToJSON(n) }
func (o *Override) As_JSON_() any       { return treeToJSON(o) }
func (o *OverrideAsList) As_JSON_() any { return treeToJSON(o) }
func (n *Node) As_JSON_() any           { return treeToJSON(n) }
