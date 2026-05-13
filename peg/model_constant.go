package peg

import asjson "github.com/neogeny/ogopego/json"

type Constant struct {
	ModelBase
	Literal string
}

func (t *Constant) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Constant) AsJSON() any               { return t.AsJSONOf(t) }

type Alert struct {
	Constant
	Level int
}

func (t *Alert) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Alert) AsJSON() any               { return t.AsJSONOf(t) }
