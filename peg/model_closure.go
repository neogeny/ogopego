package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Closure struct {
	Box
}

func (c *Closure) Parse(ctx Ctx) (trees.Tree, error) {
	return repeat(ctx, c.Exp, false)
}

func (t *Closure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Closure) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Closure) AsJSONStr() string          { return t.AsJSONStrOf(t) }
