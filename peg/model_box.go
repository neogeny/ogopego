package peg

import asjson "github.com/neogeny/ogopego/json"

type Box struct {
	ModelBase
	Exp Model
}

func (t *Box) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Box) AsJSON() any                { return t.AsJSONOf(t) }
