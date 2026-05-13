package peg

import asjson "github.com/neogeny/ogopego/json"

type Token struct {
	ModelBase
	Token string
}

func (t *Token) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Token) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Token) AsJSONStr() string          { return t.AsJSONStrOf(t) }
