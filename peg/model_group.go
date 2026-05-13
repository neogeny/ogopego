package peg

import asjson "github.com/neogeny/ogopego/json"

type Group struct {
	Box
}

func (t *Group) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Group) AsJSON() any               { return t.AsJSONOf(t) }
