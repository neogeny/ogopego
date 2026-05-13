package peg

import asjson "github.com/neogeny/ogopego/json"

type Optional struct {
	Box
}

func (t *Optional) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Optional) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Optional) AsJSONStr() string          { return t.AsJSONStrOf(t) }
