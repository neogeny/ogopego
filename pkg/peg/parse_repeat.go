// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/pkg/trees"
)

// repeat parses an expression zero or more times (or one or more if positive is true).
func repeat(ctx Ctx, exp Model, positive bool) (trees.Tree, error) {
	var items []trees.Tree

	if positive {
		ctx.CutStackPush()
		first, err := exp.Parse(ctx)
		ctx.CutStackPop()
		if err != nil {
			return nil, err
		}
		items = append(items, first)
	}

	for {
		mark := ctx.Mark()
		ctx.CutStackPush()
		result, err := exp.Parse(ctx)
		cutSeen := ctx.CutStackPop()
		if err != nil {
			if cutSeen {
				return nil, err
			}
			break
		}
		if ctx.Mark() == mark {
			return nil, ctx.Failure(
				mark,
				fmt.Errorf("closure did not consume any input, "+
					"which would lead to an infinite loop"),
			)
		}
		items = append(items, result)
	}
	return &trees.List{Items: items}, nil
}

// repeatWithSep parses an expression separated by another expression.
func repeatWithSep(
	ctx Ctx,
	exp Model,
	sep Model,
	positive bool,
	keepsep bool,
) (trees.Tree, error) {
	var items []trees.Tree

	ctx.CutStackPush()
	first, err := exp.Parse(ctx)
	ctx.CutStackPop()
	if err != nil {
		if positive {
			return nil, err
		}
		return &trees.List{Items: nil}, nil
	}
	items = append(items, first)

	for {
		mark := ctx.Mark()

		ctx.CutStackPush()
		result, err := sep.Parse(ctx)
		cutSeen := ctx.CutStackPop()

		if err != nil {
			if cutSeen {
				return nil, err
			}
			break
		}
		if keepsep {
			items = append(items, result)
		}

		ctx.CutStackPush()
		result, err = exp.Parse(ctx)
		ctx.CutStackPop()

		if err != nil {
			return nil, err
		}
		if ctx.Mark() == mark {
			return nil, ctx.Failure(
				mark,
				fmt.Errorf("closure did not consume any input, which would lead to an infinite loop"),
			)
		}
		items = append(items, result)
	}
	return &trees.List{Items: items}, nil
}
