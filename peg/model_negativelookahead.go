package peg

import (
	"fmt"
)

type NegativeLookahead struct {
	Box
}

func (n *NegativeLookahead) Parse(ctx Ctx) (Tree, error) {
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
	return NIL, nil
}

func (t *NegativeLookahead) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NegativeLookahead) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NegativeLookahead) AsJSONStr() string   { return t.AsJSONStrOf(t) }
