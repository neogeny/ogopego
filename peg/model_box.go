package peg

type Box struct {
	ModelBase
	Exp Model
}

func (t *Box) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Box) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Box) AsJSONStr() string   { return t.AsJSONStrOf(t) }
