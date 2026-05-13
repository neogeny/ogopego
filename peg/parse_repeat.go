package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

func (c *Closure) Parse(ctx Ctx) (trees.Tree, error) {
	return repeat(ctx, c.Exp, false)
}

func (p *PositiveClosure) Parse(ctx Ctx) (trees.Tree, error) {
	return repeat(ctx, p.Exp, true)
}

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
			// NOTE must match afer sep matched
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

func (j *Join) Parse(ctx Ctx) (trees.Tree, error) {
	return repeatWithSep(ctx, j.Exp, j.Sep, false, true)
}

func (p *PositiveJoin) Parse(ctx Ctx) (trees.Tree, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, true)
}

func (g *Gather) Parse(ctx Ctx) (trees.Tree, error) {
	return repeatWithSep(ctx, g.Exp, g.Sep, false, false)
}

func (g *PositiveGather) Parse(ctx Ctx) (trees.Tree, error) {
	return repeatWithSep(ctx, g.Exp, g.Sep, true, false)
}
