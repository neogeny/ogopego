package peg

import asjson "github.com/neogeny/ogopego/json"

type Lookahead struct {
	Box
}

func (t *Lookahead) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Lookahead) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Lookahead) AsJSONStr() string          { return t.AsJSONStrOf(t) }
