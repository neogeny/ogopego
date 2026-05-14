package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type NegativeLookahead struct {
	Box
}

func (n *NegativeLookahead) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, ctx.Failure(
			mark,
			fmt.Errorf(
				"negative lookahead matched:%v",
				n.Exp,
			),
		)
	}
	return trees.NIL, nil
}

func (t *NegativeLookahead) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *NegativeLookahead) AsJSON() any                { return t.AsJSONOf(t) }
func (t *NegativeLookahead) AsJSONStr() string          { return t.AsJSONStrOf(t) }
