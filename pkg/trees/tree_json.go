// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/neogeny/ogopego/pkg/asjson"
)

var _ asjson.AsJSONMixin = (*Node)(nil)

// / treeToJSONStr
func treeToJSON(t any, seen ...map[uintptr]bool) any {
	if t == nil {
		return nil
	}
	switch v := t.(type) {
	case *typeBottomTree:
		return nil
	case *treeSeq:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = treeToJSON(item, seen...)
		}
		return items
	case []any:
		items := make([]any, len(v))
		for i, item := range v {
			items[i] = treeToJSON(item, seen...)
		}
		return items
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, val := range v {
			out[k] = treeToJSON(val, seen...)
		}
		return out
	case *orderedmap.OrderedMap[string, any]:
		out := make(map[string]any, v.Len())
		for pair := v.Oldest(); pair != nil; pair = pair.Next() {
			out[pair.Key] = treeToJSON(pair.Value, seen...)
		}
		return out
	case *Node:
		child := treeToJSON(v.Tree, seen...)
		if m, ok := child.(map[string]any); ok {
			if _, has := m["__class__"]; !has {
				m["__class__"] = v.TypeName
				return m
			}
		}
		return map[string]any{"__class__": v.TypeName, "ast": child}
	default:
		return asjson.AsJSON(t, seen...)
	}
}

// As_JSON_ implementations for each concrete tree type.

func (*typeBottomTree) As_JSON_(seen map[uintptr]bool) any { return BOTTOM }
func (s *treeSeq) As_JSON_(seen map[uintptr]bool) any      { return treeToJSON(s, seen) }
func (n *Node) As_JSON_(seen map[uintptr]bool) any         { return treeToJSON(n, seen) }
