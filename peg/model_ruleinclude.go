package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type RuleInclude struct {
	ModelBase
	Name string
	Exp  Model
}

func (r *RuleInclude) Parse(ctx Ctx) (trees.Tree, error) {
	if r.Exp == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("RuleInclude %q has not been resolved", r.Name))
	}
	return r.Exp.Parse(ctx)
}

func (t *RuleInclude) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *RuleInclude) AsJSON() any                { return t.AsJSONOf(t) }
func (t *RuleInclude) AsJSONStr() string          { return t.AsJSONStrOf(t) }
