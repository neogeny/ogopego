package peg

type Lookahead struct {
	Box
}

func (l *Lookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

func (t *Lookahead) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Lookahead) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Lookahead) AsJSONStr() string   { return t.AsJSONStrOf(t) }
