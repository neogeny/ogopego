package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type NULL struct {
	ModelBase
}

func (n *NULL) Parse(ctx Ctx) (Tree, error) {
	return &trees.Nil{}, nil
}

func (t *NULL) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NULL) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NULL) AsJSONStr() string   { return t.AsJSONStrOf(t) }
