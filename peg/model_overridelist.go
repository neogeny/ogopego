package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type OverrideList struct {
	Box
}

func (o *OverrideList) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

func (t *OverrideList) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *OverrideList) AsJSON() any         { return t.AsJSONOf(t) }
func (t *OverrideList) AsJSONStr() string   { return t.AsJSONStrOf(t) }
