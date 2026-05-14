package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Synth struct {
	Box
}

func (s *Synth) Parse(ctx Ctx) (trees.Tree, error) {
	return s.Exp.Parse(ctx)
}

func (t *Synth) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Synth) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Synth) AsJSONStr() string          { return t.AsJSONStrOf(t) }
