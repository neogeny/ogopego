// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Dot matches any single rune (a dot in PEG) and returns it as text.
type Dot struct {
	ModelBase
}

// Parse implements the Model interface for Dot.
func (d *Dot) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	r, err := ctx.Dot()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &trees.Text{Value: string(r)}, nil
}
