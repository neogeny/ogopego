package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Token struct {
	ModelBase
	Token string
}

func (t *Token) Parse(ctx Ctx) (trees.Tree, error) {
	matched, err := ctx.Token(t.Token)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (t *Token) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Token) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Token) AsJSONStr() string          { return t.AsJSONStrOf(t) }
