// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

// Lookahead performs a positive lookahead: it succeeds if the nested
// expression matches without consuming input.
type Lookahead struct {
	ModelBase
	Exp Model
}

// NegativeLookahead succeeds when the nested expression does not match.
type NegativeLookahead struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for Lookahead.
func (l *Lookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	defer ctx.Reset(mark)
	ctx.EnterLookahead()
	defer ctx.LeaveLookahead()

	_, err := l.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

// Parse implements the Model interface for NegativeLookahead.
func (n *NegativeLookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	defer ctx.Reset(mark)
	ctx.EnterLookahead()
	defer ctx.LeaveLookahead()

	_, err := n.Exp.Parse(ctx)
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
