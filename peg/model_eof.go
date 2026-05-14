package peg

import (
	"fmt"
)

type EOF struct {
	ModelBase
}

func (e *EOF) Parse(ctx Ctx) (Tree, error) {
	if !ctx.Eof() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOF"),
		)
	}
	return NIL, nil
}

func (t *EOF) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EOF) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EOF) AsJSONStr() string   { return t.AsJSONStrOf(t) }
