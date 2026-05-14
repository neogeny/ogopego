package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type NamedList struct {
	Named
}

func (n *NamedList) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}

func (t *NamedList) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NamedList) AsJSON() any                { return t.AsJSONOf(t) }
func (t *NamedList) AsJSONStr() string          { return t.AsJSONStrOf(t) }
