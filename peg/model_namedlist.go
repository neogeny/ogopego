package peg

import asjson "github.com/neogeny/ogopego/json"

type NamedList struct {
	Named
}

func (t *NamedList) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NamedList) AsJSON() any                { return t.AsJSONOf(t) }
func (t *NamedList) AsJSONStr() string          { return t.AsJSONStrOf(t) }
