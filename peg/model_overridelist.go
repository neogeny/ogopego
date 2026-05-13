package peg

import asjson "github.com/neogeny/ogopego/json"

type OverrideList struct {
	Box
}

func (t *OverrideList) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *OverrideList) AsJSON() any               { return t.AsJSONOf(t) }
