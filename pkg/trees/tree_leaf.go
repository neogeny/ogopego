// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

var BOTTOM = &typeBottomTree{}

// Text is a leaf node representing plain text.
type Text struct {
	Value string
}

func (Text) tree() {}

// typeBottomTree is an internal sentinel node.
type typeBottomTree struct {
}

func (typeBottomTree) tree() {}

// TrueValue A JSON-compatible true value.
type TrueValue struct {
}

var TRUE Tree = &TrueValue{}

func (TrueValue) tree() {}

// FalseValue A JSON-compatible false value.
type FalseValue struct {
}

var FALSE Tree = &FalseValue{}

func (FalseValue) tree() {}

// NullValue A JSON-compatible null value.
type NullValue struct {
}

var NULL Tree = &NullValue{}

func (NullValue) tree() {}

// Number represents JSON-compatible a numeric literal node.
type Number struct {
	Value float64
}

func (Number) tree() {}
