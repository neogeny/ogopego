package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Fail struct {
	ModelBase
}

func (f *Fail) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Fail()
}

func (t *Fail) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Fail) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Fail) AsJSONStr() string          { return t.AsJSONStrOf(t) }
