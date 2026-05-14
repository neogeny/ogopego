// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

import asjson "github.com/neogeny/ogopego/json"

type Node struct {
	TreeBase
	TypeName string
	Tree     Tree
}

func (*Node) tree()                         {}
func (r *Node) fold(gather *treeMerge) Tree { return r }
func (r *Node) PubMap() *asjson.OrderedMap  { return r.PubMapOf(r) }
func (r *Node) AsJSON() any {
	child := r.Tree.AsJSON()
	if om, ok := child.(*asjson.OrderedMap); ok {
		if _, has := om.Get("__class__"); !has {
			out := newOM()
			out.Set("__class__", r.TypeName)
			for _, k := range om.Keys() {
				v, _ := om.Get(k)
				out.Set(k, v)
			}
			out.Set("__class__", r.TypeName)
			return out
		}
	}
	if m, ok := child.(map[string]any); ok {
		if _, has := m["__class__"]; !has {
			out := newOM()
			out.Set("__class__", r.TypeName)
			for k, v := range m {
				out.Set(k, v)
			}
			out.Set("__class__", r.TypeName)
			return out
		}
	}
	out := newOM()
	out.Set("__class__", r.TypeName)
	out.Set("ast", child)
	return out
}
func (r *Node) AsJSONStr() string { return treeJSONStr(r.AsJSON()) }
