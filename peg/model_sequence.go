package peg

import asjson "github.com/neogeny/ogopego/json"

type Sequence struct {
	ModelBase
	Sequence []Model
}

func (t *Sequence) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Sequence) AsJSON() any               { return t.AsJSONOf(t) }
