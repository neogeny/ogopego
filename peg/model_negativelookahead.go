package peg

import asjson "github.com/neogeny/ogopego/json"

type NegativeLookahead struct {
	Box
}

func (t *NegativeLookahead) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NegativeLookahead) AsJSON() any                { return t.AsJSONOf(t) }
