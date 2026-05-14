package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Lookahead struct {
	Box
}

func (l *Lookahead) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	_, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return trees.NIL, nil
}

func (t *Lookahead) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Lookahead) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Lookahead) AsJSONStr() string          { return t.AsJSONStrOf(t) }
