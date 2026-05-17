// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Fail represents a parsing failure.
type Fail struct {
	ModelBase
}

// Parse implements the Model interface for Fail.
func (f *Fail) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Fail()
}
