// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Cut represents a cut operator that prunes alternative backtracking paths.
type Cut struct {
	ModelBase
}

// Parse implements the Model interface for Cut.
func (c *Cut) Parse(ctx Ctx) (any, error) {
	ctx.Cut()
	return nil, nil
}
