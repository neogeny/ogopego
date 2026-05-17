// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

import (
	asjson "github.com/neogeny/ogopego/json"

	util "github.com/neogeny/ogopego/util"
)

// Named represents a named key/value pair folded from the parse tree.
type Named struct {
	TreeBase
	Name  string
	Value Tree
}

func (*Named) tree() {}
func (n *Named) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insert(n.Name, val)
	return val
}
func (n *Named) PubMap() *asjson.OrderedMap { return util.PubMapOf(n) }
func (n *Named) AsJSON() any {
	return map[string]any{n.Name: n.Value.AsJSON()}
}
func (n *Named) AsJSONStr() string { return treeJSONStr(n.AsJSON()) }

// NamedAsList is like Named but its values are collected as a list.
type NamedAsList struct {
	TreeBase
	Name  string
	Value Tree
}

func (*NamedAsList) tree() {}
func (n *NamedAsList) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insertAsList(n.Name, val)
	return val
}
func (n *NamedAsList) PubMap() *asjson.OrderedMap { return util.PubMapOf(n) }
func (n *NamedAsList) AsJSON() any {
	return map[string]any{n.Name: n.Value.AsJSON()}
}
func (n *NamedAsList) AsJSONStr() string { return treeJSONStr(n.AsJSON()) }

// Override indicates that the contained value should override other values
// when folding into the result.
type Override struct {
	TreeBase
	Value Tree
}

func (*Override) tree() {}
func (o *Override) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendTree(gather.Root, val)
	return val
}
func (o *Override) PubMap() *asjson.OrderedMap { return util.PubMapOf(o) }
func (o *Override) AsJSON() any                { return o.Value.AsJSON() }
func (o *Override) AsJSONStr() string          { return treeJSONStr(o.AsJSON()) }

// OverrideAsList is a list-form override variant.
type OverrideAsList struct {
	TreeBase
	Value Tree
}

func (*OverrideAsList) tree() {}
func (o *OverrideAsList) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendAsList(gather.Root, val)
	return val
}
func (o *OverrideAsList) PubMap() *asjson.OrderedMap { return util.PubMapOf(o) }
func (o *OverrideAsList) AsJSON() any                { return o.Value.AsJSON() }
func (o *OverrideAsList) AsJSONStr() string          { return treeJSONStr(o.AsJSON()) }
