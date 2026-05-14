package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Named struct {
	NamedBox
}

func (n *Named) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

func (t *Named) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Named) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Named) AsJSONStr() string   { return t.AsJSONStrOf(t) }
