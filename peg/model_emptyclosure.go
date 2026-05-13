package peg

import asjson "github.com/neogeny/ogopego/json"

type EmptyClosure struct {
	ModelBase
}

func (t *EmptyClosure) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EmptyClosure) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EmptyClosure) AsJSONStr() string          { return t.AsJSONStrOf(t) }
