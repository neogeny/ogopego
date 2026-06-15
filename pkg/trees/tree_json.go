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
	case *typeBottomTree:
		return nil
	case *TreeSeq:
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

func (*typeBottomTree) As_JSON_() any { return BOTTOM }
func (s *TreeSeq) As_JSON_() any          { return treeToJSON(s) }
func (n *Node) As_JSON_() any         { return treeToJSON(n) }
