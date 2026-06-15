// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

var NIL = &Nil{}
var BOTTOM = &Bottom{}

// Text is a leaf node representing plain text.
type Text struct {
	Value string
}

func (Text) tree()                         {}
func (t *Text) fold(gather *treeMerge) any { return t }

// Bool represents a boolean node.
type Bool struct {
	Value bool
}

func (Bool) tree()                         {}
func (b *Bool) fold(gather *treeMerge) any { return b }

// Nil is the nil sentinel node used to represent empty results.
type Nil struct {
}

func (Nil) tree()                         {}
func (n *Nil) fold(gather *treeMerge) any { return n }

// Bottom is an internal sentinel node.
type Bottom struct {
}

func (Bottom) tree()                         {}
func (b *Bottom) fold(gather *treeMerge) any { return b }

// TrueValue A JSON-compatible true value.
type TrueValue struct {
}

var TRUE Tree = &TrueValue{}

func (TrueValue) tree()                         {}
func (v *TrueValue) fold(gather *treeMerge) any { return v }

// FalseValue A JSON-compatible false value.
type FalseValue struct {
}

var FALSE Tree = &FalseValue{}

func (FalseValue) tree()                         {}
func (v *FalseValue) fold(gather *treeMerge) any { return v }

// NullValue A JSON-compatible null value.
type NullValue struct {
}

var NULL Tree = &NullValue{}

func (NullValue) tree()                         {}
func (v *NullValue) fold(gather *treeMerge) any { return v }

// Number represents JSON-compatible a numeric literal node.
type Number struct {
	Value float64
}

func (Number) tree()                         {}
func (n *Number) fold(gather *treeMerge) any { return n }
