package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type EmptyClosure struct {
	ModelBase
}

func (e *EmptyClosure) Parse(ctx Ctx) (trees.Tree, error) {
	return &trees.List{Items: nil}, nil
}

func (t *EmptyClosure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EmptyClosure) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EmptyClosure) AsJSONStr() string          { return t.AsJSONStrOf(t) }
