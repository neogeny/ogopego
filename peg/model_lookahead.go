// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// Lookahead performs a positive lookahead: it succeeds if the nested
// expression matches without consuming input.
type Lookahead struct {
	Box
}

// NegativeLookahead succeeds when the nested expression does not match.
type NegativeLookahead struct {
	Box
}

// Parse implements the Model interface for Lookahead.
func (l *Lookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

// Parse implements the Model interface for NegativeLookahead.
func (n *NegativeLookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, ctx.Failure(
			mark,
			fmt.Errorf(
				"negative lookahead matched:%v",
				n.Exp,
			),
		)
	}
	return NIL, nil
}
