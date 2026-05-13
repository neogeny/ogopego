package peg

import asjson "github.com/neogeny/ogopego/json"

type RuleInclude struct {
	ModelBase
	Name string
	Exp  Model
}

func (t *RuleInclude) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *RuleInclude) AsJSON() any                { return t.AsJSONOf(t) }
