package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

type Token struct {
	ModelBase
	Token string
}

func (t *Token) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	if !ctx.MatchToken(t.Token) {
		ctx.Reset(mark)
		return nil, ctx.Failure(mark, fmt.Errorf("expected: %q", t.Token))
	}
	return &trees.Text{Value: t.Token}, nil
}

func (t *Token) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Token) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Token) AsJSONStr() string   { return t.AsJSONStrOf(t) }
