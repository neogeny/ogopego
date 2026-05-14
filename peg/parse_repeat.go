package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

func repeat(ctx Ctx, exp Model, positive bool) (trees.Tree, error) {
	var items []trees.Tree

	if positive {
		first, err := exp.Parse(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, first)
	}

	for {
		mark := ctx.Mark()
		result, err := exp.Parse(ctx)
		if err != nil {
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

func repeatWithSep(
	ctx Ctx,
	exp Model,
	sep Model,
	positive bool,
	keepsep bool,
) (trees.Tree, error) {
	var items []trees.Tree

	first, err := exp.Parse(ctx)
	if err != nil {
		if positive {
			return nil, err
		}
		return &trees.List{Items: nil}, nil
	}
	items = append(items, first)

	for {
		mark := ctx.Mark()
		result, err := sep.Parse(ctx)
		if err != nil {
			break
		}
		if keepsep {
			items = append(items, result)
		}

		result, err = exp.Parse(ctx)
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
