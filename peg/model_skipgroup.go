package peg

import asjson "github.com/neogeny/ogopego/json"

type SkipGroup struct {
	Box
}

func (t *SkipGroup) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *SkipGroup) AsJSON() any                { return t.AsJSONOf(t) }
func (t *SkipGroup) AsJSONStr() string          { return t.AsJSONStrOf(t) }
