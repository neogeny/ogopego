package peg

import asjson "github.com/neogeny/ogopego/json"

type Void struct {
	ModelBase
}

func (t *Void) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Void) AsJSON() any               { return t.AsJSONOf(t) }
