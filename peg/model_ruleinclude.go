package peg

import (
	"fmt"
)

type RuleInclude struct {
	ModelBase
	Name string
	Exp  Model
}

func (r *RuleInclude) Parse(ctx Ctx) (Tree, error) {
	if r.Exp == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("RuleInclude %q has not been resolved", r.Name))
	}
	return r.Exp.Parse(ctx)
}

func (t *RuleInclude) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *RuleInclude) AsJSON() any         { return t.AsJSONOf(t) }
func (t *RuleInclude) AsJSONStr() string   { return t.AsJSONStrOf(t) }
