package peg

import asjson "github.com/neogeny/ogopego/json"

type Override struct {
	Box
}

func (t *Override) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Override) AsJSON() any               { return t.AsJSONOf(t) }
