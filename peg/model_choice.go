package peg

import (
	"errors"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Choice struct {
	ModelBase
	Options []*Option
}

func (c *Choice) Parse(ctx Ctx) (trees.Tree, error) {
	startMark := ctx.Mark()
	var lastErr error
	for _, opt := range c.Options {
		mark := ctx.Mark()
		result, err := opt.Parse(ctx)
		if err == nil {
			return result, nil
		}
		ctx.Reset(mark)
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no option matched")
	}
	ctx.Reset(startMark)
	return nil, lastErr
}

func (o *Optional) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return &trees.Nil{}, nil
	}
	return result, nil
}

func (c *Choice) PubMap() *asjson.OrderedMap { return c.PubMapOf(c) }
func (c *Choice) AsJSON() any                { return c.AsJSONOf(c) }

func (o *Option) PubMap() *asjson.OrderedMap { return o.PubMapOf(o) }
func (o *Option) AsJSON() any                { return o.AsJSONOf(o) }
