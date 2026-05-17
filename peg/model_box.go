// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Box is a model wrapper that contains a nested expression (Exp).
type Box struct {
	ModelBase
	Exp Model
}

// NamedBox is a Box that carries a name for the nested expression.
type NamedBox struct {
	Box
	Name string
}
