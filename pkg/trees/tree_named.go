// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// Named represents a named key/value pair folded from the parse tree.
type Named struct {
	Name  string
	Value any
}

func (Named) tree() {}

// NamedAsList is like Named but its values are collected as a list.
type NamedAsList struct {
	Name  string
	Value any
}

func (NamedAsList) tree() {}

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
