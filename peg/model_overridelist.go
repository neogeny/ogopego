package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type OverrideList struct {
	Box
}

func (o *OverrideList) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

func (t *OverrideList) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *OverrideList) AsJSON() any                { return t.AsJSONOf(t) }
func (t *OverrideList) AsJSONStr() string          { return t.AsJSONStrOf(t) }
