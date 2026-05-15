package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Pattern struct {
	ModelBase
	Pattern string
}

func (p *Pattern) Parse(ctx Ctx) (Tree, error) {
	matched, err := ctx.MatchPattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (t *Pattern) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Pattern) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Pattern) AsJSONStr() string   { return t.AsJSONStrOf(t) }
