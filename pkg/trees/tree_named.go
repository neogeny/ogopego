// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// Named represents a named key/value pair folded from the parse tree.
type Named struct {
	Name  string
	Value any
}

func (Named) tree() {}
func (n *Named) fold(gather *treeMerge) any {
	val := fold(gather, n.Value)
	gather.insert(n.Name, val)
	return val
}

// NamedAsList is like Named but its values are collected as a list.
type NamedAsList struct {
	Name  string
	Value any
}

func (NamedAsList) tree() {}
func (n *NamedAsList) fold(gather *treeMerge) any {
	val := fold(gather, n.Value)
	gather.insertAsList(n.Name, val)
	return val
}

// Override indicates that the contained value should override other values
// when folding into the result.
type Override struct {
	Value any
}

func (Override) tree() {}
func (o *Override) fold(gather *treeMerge) any {
	val := fold(gather, o.Value)
	gather.Root = appendTree(gather.Root, val)
	return val
}

// OverrideAsList is a list-form override variant.
type OverrideAsList struct {
	Value any
}

func (OverrideAsList) tree() {}
func (o *OverrideAsList) fold(gather *treeMerge) any {
	val := fold(gather, o.Value)
	gather.Root = appendAsSeq(gather.Root, val)
	return val
}
