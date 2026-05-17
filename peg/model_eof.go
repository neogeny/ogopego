// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// EOF matches the end of input.
type EOF struct {
	ModelBase
}

// Parse implements the Model interface for EOF.
func (e *EOF) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	ctx.NextToken()
	if !ctx.Eof() {
		ctx.Reset(mark)
		return nil, ctx.Failure(
			mark,
			fmt.Errorf("expected EOF"),
		)
	}
	return NIL, nil
}
