package peg

import asjson "github.com/neogeny/ogopego/json"

type Cut struct {
	ModelBase
}

func (t *Cut) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Cut) AsJSON() any               { return t.AsJSONOf(t) }
