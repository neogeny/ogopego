// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// Text is a leaf node representing plain text.
type Text struct {
	Value string
}

func (Text) tree() {}

// Number represents JSON-compatible a numeric literal node.
type Number struct {
	Value float64
}

func (Number) tree() {}
