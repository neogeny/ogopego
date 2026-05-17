// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Cut represents a cut operator that prunes alternative backtracking paths.
type Cut struct {
	ModelBase
}

// Parse implements the Model interface for Cut.
func (c *Cut) Parse(ctx Ctx) (Tree, error) {
	ctx.Cut()
	t := &trees.Nil{}
	return t, nil
}
