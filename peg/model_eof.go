package peg

import (
	"fmt"
)

type EOF struct {
	ModelBase
}

func (e *EOF) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	ctx.NextToken()
	if !ctx.Eof() {
		ctx.Reset(mark)
		return nil, ctx.Failure(
			mark,
			fmt.Errorf("expected EOF"),
		)
	}
	return NIL, nil
}

func (t *EOF) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EOF) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EOF) AsJSONStr() string   { return t.AsJSONStrOf(t) }
