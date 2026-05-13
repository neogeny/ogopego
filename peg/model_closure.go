package peg

import asjson "github.com/neogeny/ogopego/json"

type Closure struct {
	Box
}

func (t *Closure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Closure) AsJSON() any                { return t.AsJSONOf(t) }
