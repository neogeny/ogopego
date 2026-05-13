package peg

import asjson "github.com/neogeny/ogopego/json"

type NULL struct {
	ModelBase
}

func (t *NULL) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NULL) AsJSON() any                { return t.AsJSONOf(t) }
func (t *NULL) AsJSONStr() string          { return t.AsJSONStrOf(t) }
