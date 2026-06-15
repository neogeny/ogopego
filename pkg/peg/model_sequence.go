// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/pkg/trees"
)

// Sequence represents an ordered sequence of model elements.
type Sequence struct {
	ModelBase
	Sequence []Model
}

// Parse implements the Model interface for Sequence.
func (s *Sequence) Parse(ctx Ctx) (any, error) {
	mark := ctx.Mark()
	var out any
	for _, el := range s.Sequence {
		tree, err := el.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			return nil, err
		}
		if tree == nil {
			continue
		}
		out = trees.MergeTrees(out, tree)
	}
	return out, nil
}
