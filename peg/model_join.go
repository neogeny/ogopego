package peg

import asjson "github.com/neogeny/ogopego/json"

type Join struct {
	Box
	Sep Model
}

func (t *Join) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Join) AsJSON() any               { return t.AsJSONOf(t) }

type PositiveJoin struct {
	Join
}

func (t *PositiveJoin) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *PositiveJoin) AsJSON() any               { return t.AsJSONOf(t) }

type Gather struct {
	Join
}

func (t *Gather) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Gather) AsJSON() any               { return t.AsJSONOf(t) }

type PositiveGather struct {
	Gather
}

func (t *PositiveGather) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *PositiveGather) AsJSON() any               { return t.AsJSONOf(t) }
