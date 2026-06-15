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

// Override indicates that the contained value should override other values
// when folding into the result.
type Override struct {
	Value any
}

func (Override) tree() {}

// OverrideAsList is a list-form override variant.
type OverrideAsList struct {
	Value any
}

func (OverrideAsList) tree() {}
