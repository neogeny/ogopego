// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

import asjson "github.com/neogeny/ogopego/json"

// Seq represents a sequence node whose items are merged when folding.
type Seq struct {
	TreeBase
	Items []Tree
}

func (*Seq) tree() {}
func (s *Seq) fold(gather *treeMerge) Tree {
	var out Tree = &Nil{}
	for _, item := range s.Items {
		out = merge(out, item.fold(gather))
	}
	return out
}
func (s *Seq) PubMap() *asjson.OrderedMap { return s.PubMapOf(s) }
func (s *Seq) AsJSON() any {
	items := make([]any, len(s.Items))
	for i, item := range s.Items {
		items[i] = item.AsJSON()
	}
	return items
}
func (s *Seq) AsJSONStr() string { return treeJSONStr(s.AsJSON()) }

// List represents a closed list node produced after folding sequences.
type List struct {
	TreeBase
	Items []Tree
}

func (*List) tree() {}
func (l *List) fold(gather *treeMerge) Tree {
	items := make([]Tree, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.fold(gather)
	}
	return &List{Items: items}
}
func (l *List) PubMap() *asjson.OrderedMap { return l.PubMapOf(l) }
func (l *List) AsJSON() any {
	items := make([]any, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.AsJSON()
	}
	return items
}
func (l *List) AsJSONStr() string { return treeJSONStr(l.AsJSON()) }
