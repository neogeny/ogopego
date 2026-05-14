package peg

type Closure struct {
	Box
}

func (c *Closure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, c.Exp, false)
}

func (t *Closure) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Closure) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Closure) AsJSONStr() string   { return t.AsJSONStrOf(t) }
