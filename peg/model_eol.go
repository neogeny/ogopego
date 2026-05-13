package peg

import asjson "github.com/neogeny/ogopego/json"

type EOL struct {
	ModelBase
}

func (t *EOL) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOL) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EOL) AsJSONStr() string          { return t.AsJSONStrOf(t) }
