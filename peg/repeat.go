package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

func errNotImplemented(name string) error {
	return fmt.Errorf("%s not yet implemented", name)
}

func (c *Closure) Parse(ctx Ctx) (trees.Tree, error) {
	var items []trees.Tree
	for {
		mark := ctx.Mark()
		result, err := c.Exp.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			break
		}
		if ctx.Mark() == mark {
			break
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}

func (p *PositiveClosure) Parse(ctx Ctx) (trees.Tree, error) {
	startMark := ctx.Mark()
	first, err := p.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(startMark)
		return nil, err
	}
	var items []trees.Tree
	if _, ok := first.(*trees.Nil); !ok {
		items = append(items, first)
	}
	for {
		mark := ctx.Mark()
		result, err := p.Exp.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			break
		}
		if ctx.Mark() == mark {
			break
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}

func (j *Join) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, errNotImplemented("Join")
}

func (p *PositiveJoin) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, errNotImplemented("PositiveJoin")
}

func (g *Gather) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, errNotImplemented("Gather")
}

func (p *PositiveGather) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, errNotImplemented("PositiveGather")
}
