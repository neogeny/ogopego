// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

var NIL = &Nil{}
var BOTTOM = &Bottom{}

// Text is a leaf node representing plain text.
type Text struct {
	TreeBase
	Value string
}

func (*Text) tree()                         {}
func (t *Text) fold(gather *treeMerge) Tree { return t }

// Number represents a numeric literal node.
type Number struct {
	TreeBase
	Value float64
}

func (*Number) tree()                         {}
func (n *Number) fold(gather *treeMerge) Tree { return n }

// Bool represents a boolean node.
type Bool struct {
	TreeBase
	Value bool
}

func (*Bool) tree()                         {}
func (b *Bool) fold(gather *treeMerge) Tree { return b }

// Nil is the nil sentinel node used to represent empty results.
type Nil struct {
	TreeBase
}

func (*Nil) tree()                         {}
func (n *Nil) fold(gather *treeMerge) Tree { return n }

// Bottom is an internal sentinel node.
type Bottom struct {
	TreeBase
}

func (*Bottom) tree()                         {}
func (b *Bottom) fold(gather *treeMerge) Tree { return b }
