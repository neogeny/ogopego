package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/tree"
)

type Call struct {
	ModelBase
	Name   string
	Target *Rule
}

func (c *Call) Parse(ctx Ctx) (tree.Tree, error) {
	if c.Target == nil {
		return nil, fmt.Errorf("call to %q has not been linked", c.Name)
	}
	mark := ctx.Mark()
	result, err := c.Target.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return result, nil
}
