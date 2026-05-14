package peg

type Join struct {
	Box
	Sep Model
}

type PositiveJoin struct {
	Join
}

type Gather struct {
	Join
}

type PositiveGather struct {
	Gather
}

func (j *Join) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, j.Exp, j.Sep, false, true)
}

func (t *Join) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Join) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Join) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (p *PositiveJoin) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, true)
}

func (t *PositiveJoin) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *PositiveJoin) AsJSON() any         { return t.AsJSONOf(t) }
func (t *PositiveJoin) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (g *Gather) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, g.Exp, g.Sep, false, false)
}

func (t *Gather) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Gather) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Gather) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (p *PositiveGather) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, false)
}

func (t *PositiveGather) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *PositiveGather) AsJSON() any         { return t.AsJSONOf(t) }
func (t *PositiveGather) AsJSONStr() string   { return t.AsJSONStrOf(t) }
