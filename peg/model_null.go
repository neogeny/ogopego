package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type NULL struct {
	ModelBase
}

func (n *NULL) Parse(ctx Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (t *NULL) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NULL) AsJSON() any                { return t.AsJSONOf(t) }
func (t *NULL) AsJSONStr() string          { return t.AsJSONStrOf(t) }
