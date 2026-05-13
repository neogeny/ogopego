package peg

import asjson "github.com/neogeny/ogopego/json"

type Pattern struct {
	ModelBase
	Pattern string
}

func (t *Pattern) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Pattern) AsJSON() any               { return t.AsJSONOf(t) }
