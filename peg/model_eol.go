package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type EOL struct {
	ModelBase
}

func (e *EOL) Parse(ctx Ctx) (trees.Tree, error) {
	if !ctx.MatchEOL() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOL"),
		)
	}
	return &trees.Nil{}, nil
}

func (t *EOL) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOL) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EOL) AsJSONStr() string          { return t.AsJSONStrOf(t) }
