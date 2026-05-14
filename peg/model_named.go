package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Named struct {
	NamedBox
}

func (n *Named) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

func (t *Named) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Named) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Named) AsJSONStr() string          { return t.AsJSONStrOf(t) }
