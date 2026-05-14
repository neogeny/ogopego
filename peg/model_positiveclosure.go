package peg

type PositiveClosure struct {
	Closure
}

func (p *PositiveClosure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, p.Exp, true)
}

func (t *PositiveClosure) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *PositiveClosure) AsJSON() any         { return t.AsJSONOf(t) }
func (t *PositiveClosure) AsJSONStr() string   { return t.AsJSONStrOf(t) }
