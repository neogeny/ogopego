package peg

import asjson "github.com/neogeny/ogopego/json"

type EOF struct {
	ModelBase
}

func (t *EOF) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOF) AsJSON() any                { return t.AsJSONOf(t) }
