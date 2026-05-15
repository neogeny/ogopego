// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

import asjson "github.com/neogeny/ogopego/json"

var NIL = &Nil{}
var BOTTOM = &Bottom{}

type Text struct {
	TreeBase
	Value string
}

func (*Text) tree()                         {}
func (t *Text) fold(gather *treeMerge) Tree { return t }
func (t *Text) PubMap() *asjson.OrderedMap  { return t.PubMapOf(t) }
func (t *Text) AsJSON() any                 { return t.Value }
func (t *Text) AsJSONStr() string           { return treeJSONStr(t.AsJSON()) }

type Number struct {
	TreeBase
	Value float64
}

func (*Number) tree()                         {}
func (n *Number) fold(gather *treeMerge) Tree { return n }
func (n *Number) PubMap() *asjson.OrderedMap  { return n.PubMapOf(n) }
func (n *Number) AsJSON() any                 { return n.Value }
func (n *Number) AsJSONStr() string           { return treeJSONStr(n.AsJSON()) }

type Bool struct {
	TreeBase
	Value bool
}

func (*Bool) tree()                         {}
func (b *Bool) fold(gather *treeMerge) Tree { return b }
func (b *Bool) PubMap() *asjson.OrderedMap  { return b.PubMapOf(b) }
func (b *Bool) AsJSON() any                 { return b.Value }
func (b *Bool) AsJSONStr() string           { return treeJSONStr(b.AsJSON()) }

type Nil struct {
	TreeBase
}

func (*Nil) tree()                         {}
func (n *Nil) fold(gather *treeMerge) Tree { return n }
func (n *Nil) PubMap() *asjson.OrderedMap  { return n.PubMapOf(n) }
func (n *Nil) AsJSON() any                 { return nil }
func (n *Nil) AsJSONStr() string           { return treeJSONStr(n.AsJSON()) }

type Bottom struct {
	TreeBase
}

func (*Bottom) tree()                         {}
func (b *Bottom) fold(gather *treeMerge) Tree { return b }
func (b *Bottom) PubMap() *asjson.OrderedMap  { return b.PubMapOf(b) }
func (b *Bottom) AsJSON() any                 { return nil }
func (b *Bottom) AsJSONStr() string           { return treeJSONStr(b.AsJSON()) }
