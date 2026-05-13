package peg

import asjson "github.com/neogeny/ogopego/json"

type NamedBox struct {
	Box
	Name string
}

func (t *NamedBox) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NamedBox) AsJSON() any                { return t.AsJSONOf(t) }
