package peg

type NamedBox struct {
	Box
	Name string
}

func (t *NamedBox) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NamedBox) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NamedBox) AsJSONStr() string   { return t.AsJSONStrOf(t) }
