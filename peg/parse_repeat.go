package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/tree"
)

func errNotImplemented(name string) error {
	return fmt.Errorf("%s not yet implemented", name)
}

func (c *Closure) Parse(ctx Ctx) (tree.Tree, error) {
	var items []tree.Tree
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
		if _, ok := result.(*tree.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &tree.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &tree.Seq{Items: items}, nil
	}
}

func (p *PositiveClosure) Parse(ctx Ctx) (tree.Tree, error) {
	startMark := ctx.Mark()
	first, err := p.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(startMark)
		return nil, err
	}
	var items []tree.Tree
	if _, ok := first.(*tree.Nil); !ok {
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
		if _, ok := result.(*tree.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &tree.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &tree.Seq{Items: items}, nil
	}
}

func (j *Join) Parse(ctx Ctx) (tree.Tree, error) {
	return nil, errNotImplemented("Join")
}

func (p *PositiveJoin) Parse(ctx Ctx) (tree.Tree, error) {
	return nil, errNotImplemented("PositiveJoin")
}

func (g *Gather) Parse(ctx Ctx) (tree.Tree, error) {
	return nil, errNotImplemented("Gather")
}

func (p *PositiveGather) Parse(ctx Ctx) (tree.Tree, error) {
	return nil, errNotImplemented("PositiveGather")
}
