package peg

import asjson "github.com/neogeny/ogopego/json"

type PositiveClosure struct {
	Closure
}

func (t *PositiveClosure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *PositiveClosure) AsJSON() any               { return t.AsJSONOf(t) }
