// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

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
