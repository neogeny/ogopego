package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type EOF struct {
	ModelBase
}

func (e *EOF) Parse(ctx Ctx) (trees.Tree, error) {
	if !ctx.Eof() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOF"),
		)
	}
	return trees.NIL, nil
}

func (t *EOF) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOF) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EOF) AsJSONStr() string          { return t.AsJSONStrOf(t) }
