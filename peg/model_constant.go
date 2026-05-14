package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Constant struct {
	ModelBase
	Literal string
}

func (c *Constant) Parse(ctx Ctx) (trees.Tree, error) {
	return ctx.Constant(c.Literal)
}

func (t *Constant) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Constant) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Constant) AsJSONStr() string          { return t.AsJSONStrOf(t) }

type Alert struct {
	Constant
	Level int
}

func (a *Alert) Parse(ctx Ctx) (trees.Tree, error) {
	return ctx.Constant(a.Literal)
}

func (t *Alert) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Alert) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Alert) AsJSONStr() string          { return t.AsJSONStrOf(t) }
