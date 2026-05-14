package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Dot struct {
	ModelBase
}

func (d *Dot) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	r, err := ctx.Dot()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &trees.Text{Value: string(r)}, nil
}

func (t *Dot) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Dot) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Dot) AsJSONStr() string   { return t.AsJSONStrOf(t) }
