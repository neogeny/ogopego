package peg

import asjson "github.com/neogeny/ogopego/json"

type Fail struct {
	ModelBase
}

func (t *Fail) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Fail) AsJSON() any                { return t.AsJSONOf(t) }
