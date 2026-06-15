// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Dot matches any single rune (a dot in PEG) and returns it as text.
type Dot struct {
	ModelBase
}

// Parse implements the Model interface for Dot.
func (d *Dot) Parse(ctx Ctx) (any, error) {
	mark := ctx.Mark()
	r, err := ctx.Dot()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return string(r), nil
}
