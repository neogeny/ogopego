package peg

import asjson "github.com/neogeny/ogopego/json"

type Comment struct {
	ModelBase
	Comment string
}

func (t *Comment) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Comment) AsJSON() any                { return t.AsJSONOf(t) }

type EOLComment struct {
	Comment
}

func (t *EOLComment) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOLComment) AsJSON() any                { return t.AsJSONOf(t) }
