package peg

import asjson "github.com/neogeny/ogopego/json"

type SkipTo struct {
	Box
}

func (t *SkipTo) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *SkipTo) AsJSON() any                { return t.AsJSONOf(t) }
func (t *SkipTo) AsJSONStr() string          { return t.AsJSONStrOf(t) }
