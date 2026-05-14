package peg

type SkipGroup struct {
	Box
}

func (s *SkipGroup) Parse(ctx Ctx) (Tree, error) {
	_, err := s.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

func (t *SkipGroup) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *SkipGroup) AsJSON() any         { return t.AsJSONOf(t) }
func (t *SkipGroup) AsJSONStr() string   { return t.AsJSONStrOf(t) }
