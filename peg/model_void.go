package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Void struct {
	ModelBase
}

func (v *Void) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	err := ctx.Void()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return trees.NIL, nil
}

func (t *Void) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Void) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Void) AsJSONStr() string          { return t.AsJSONStrOf(t) }
