package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Pattern struct {
	ModelBase
	Pattern string
}

func (p *Pattern) Parse(ctx Ctx) (trees.Tree, error) {
	matched, err := ctx.Pattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (t *Pattern) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Pattern) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Pattern) AsJSONStr() string          { return t.AsJSONStrOf(t) }
