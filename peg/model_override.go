package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Override struct {
	Box
}

func (o *Override) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

func (t *Override) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Override) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Override) AsJSONStr() string          { return t.AsJSONStrOf(t) }
