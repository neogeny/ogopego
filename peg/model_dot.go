package peg

import asjson "github.com/neogeny/ogopego/json"

type Dot struct {
	ModelBase
}

func (t *Dot) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Dot) AsJSON() any               { return t.AsJSONOf(t) }
