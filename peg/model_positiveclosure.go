package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type PositiveClosure struct {
	Closure
}

func (p *PositiveClosure) Parse(ctx Ctx) (trees.Tree, error) {
	return repeat(ctx, p.Exp, true)
}

func (t *PositiveClosure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *PositiveClosure) AsJSON() any                { return t.AsJSONOf(t) }
func (t *PositiveClosure) AsJSONStr() string          { return t.AsJSONStrOf(t) }
