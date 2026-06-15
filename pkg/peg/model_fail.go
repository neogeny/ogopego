// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Fail represents a parsing failure.
type Fail struct {
	ModelBase
}

// Parse implements the Model interface for Fail.
func (f *Fail) Parse(ctx Ctx) (any, error) {
	return nil, ctx.Fail()
}
