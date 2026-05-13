package peg

import asjson "github.com/neogeny/ogopego/json"

type Named struct {
	NamedBox
}

func (t *Named) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Named) AsJSON() any                { return t.AsJSONOf(t) }
