package peg

import asjson "github.com/neogeny/ogopego/json"

type Synth struct {
	Box
}

func (t *Synth) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Synth) AsJSON() any                { return t.AsJSONOf(t) }
