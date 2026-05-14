package peg

import (
	"errors"

	"github.com/neogeny/ogopego/context"
)

type Option struct {
	Box
}

type Choice struct {
	ModelBase
	Options []*Option
}

func (o *Option) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Choice) Parse(ctx Ctx) (Tree, error) {
	startMark := ctx.Mark()
	var lastErr error
	for _, opt := range c.Options {
		mark := ctx.Mark()
		result, err := opt.Parse(ctx)
		if err == nil {
			return result, nil
		}
		ctx.Reset(mark)
		if context.TakeCut(err) {
			return nil, err
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = ctx.Failure(startMark, errors.New("no option matched"))
	}
	ctx.Reset(startMark)
	return nil, lastErr
}

func (c *Choice) PubMap() *OrderedMap { return c.PubMapOf(c) }
func (c *Choice) AsJSON() any         { return c.AsJSONOf(c) }
func (c *Choice) AsJSONStr() string   { return c.AsJSONStrOf(c) }

func (o *Option) PubMap() *OrderedMap { return o.PubMapOf(o) }
func (o *Option) AsJSON() any         { return o.AsJSONOf(o) }
func (o *Option) AsJSONStr() string   { return o.AsJSONStrOf(o) }
