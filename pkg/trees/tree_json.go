// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"maps"

	"github.com/neogeny/ogopego/pkg/asjson"
)

var _ asjson.AsJSONMixin = (*Node)(nil)

// / treeToJSONStr
func treeToJSON(t any) any {
	if t == nil {
		return nil
	}
	switch v := t.(type) {
	case *Text:
		return v.Value
	case *Number:
		return v.Value
	case *typeBottomTree:
		return nil
	case *Seq:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = treeToJSON(item)
		}
		return items
	case []any:
		items := make([]any, len(v))
		for i, item := range v {
			items[i] = treeToJSON(item)
		}
		return items
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, val := range v {
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
func (*typeBottomTree) As_JSON_() any   { return nil }
func (s *Seq) As_JSON_() any            { return treeToJSON(s) }
func (n *Named) As_JSON_() any          { return treeToJSON(n) }
func (n *NamedAsList) As_JSON_() any    { return treeToJSON(n) }
func (o *Override) As_JSON_() any       { return treeToJSON(o) }
func (o *OverrideAsList) As_JSON_() any { return treeToJSON(o) }
func (n *Node) As_JSON_() any           { return treeToJSON(n) }
