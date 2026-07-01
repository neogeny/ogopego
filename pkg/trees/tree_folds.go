// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"unicode"
	"unique"
)

type namedTreeT struct {
	NameHandle unique.Handle[string]
	Tree       any
}

func (*namedTreeT) isTree()                              {}
func (t *namedTreeT) As_JSON_(seen map[uintptr]bool) any { return treeToJSON(t, seen) }

type namedTreeSeqT struct {
	NameHandle unique.Handle[string]
	Tree       any
}

func (*namedTreeSeqT) isTree()                              {}
func (t *namedTreeSeqT) As_JSON_(seen map[uintptr]bool) any { return treeToJSON(t, seen) }

type ovrTreeT struct {
	Tree any
}

func (*ovrTreeT) isTree()                              {}
func (t *ovrTreeT) As_JSON_(seen map[uintptr]bool) any { return treeToJSON(t, seen) }

type overrideTreeSeqT struct {
	Tree any
}

func (*overrideTreeSeqT) isTree()                              {}
func (t *overrideTreeSeqT) As_JSON_(seen map[uintptr]bool) any { return treeToJSON(t, seen) }

func validateUserKeyName(name string) {
	for _, char := range name {
		if !unicode.IsLetter(char) && char != '_' {
			panic("invalid name: " + name)
		}
	}
}

// NamedTree represents a named key/value pair folded from the parse tree.
func NamedTree(name string, tree any) *namedTreeT {
	validateUserKeyName(name)
	return &namedTreeT{NameHandle: unique.Make(name), Tree: tree}
}

// NamedTreeSeq is like Named but its values are collected as a list.
func NamedTreeSeq(name string, tree any) *namedTreeSeqT {
	validateUserKeyName(name)
	return &namedTreeSeqT{NameHandle: unique.Make(name), Tree: tree}
}

func (t *namedTreeT) Name() string {
	return t.NameHandle.Value()
}

func (named *namedTreeSeqT) Name() string {
	return named.NameHandle.Value()
}

// OverrideTree indicates that the contained value should override other values
// when folding into the result.
func OverrideTree(tree any) *ovrTreeT {
	return &ovrTreeT{Tree: tree}
}

// OverrideTreeSeq is a list-form override variant.
func OverrideTreeSeq(tree any) *overrideTreeSeqT {
	return &overrideTreeSeqT{Tree: tree}
}
