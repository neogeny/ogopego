// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// treeNamed represents a named key/value pair folded from the parse tree.
type treeNamed struct {
	Name  string
	Value any
}

func (treeNamed) tree() {}

// treeNamedAsList is like Named but its values are collected as a list.
type treeNamedAsList struct {
	Name  string
	Value any
}

func (treeNamedAsList) tree() {}

// treeOverride indicates that the contained value should override other values
// when folding into the result.
type treeOverride struct {
	Value any
}

func (treeOverride) tree() {}

// treeOverrideAsList is a list-form override variant.
type treeOverrideAsList struct {
	Value any
}

func (treeOverrideAsList) tree() {}
