package peg

import (
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Optional struct {
	Box
}

func (o *Optional) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		if context.TakeCut(err) {
			return nil, err
		}
		return &trees.Nil{}, nil
	}
	return result, nil
}

func (t *Optional) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Optional) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Optional) AsJSONStr() string   { return t.AsJSONStrOf(t) }
