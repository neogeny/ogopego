package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Cut struct {
	ModelBase
}

func (c *Cut) Parse(ctx Ctx) (trees.Tree, error) {
	ctx.Tracer().TraceCut(ctx)
	t := &trees.Nil{}
	t.OrCutSeen(true)
	return t, nil
}

func (t *Cut) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Cut) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Cut) AsJSONStr() string          { return t.AsJSONStrOf(t) }
