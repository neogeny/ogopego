package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type SkipGroup struct {
	Box
}

func (s *SkipGroup) Parse(ctx Ctx) (Tree, error) {
	_, err := s.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return trees.NIL, nil
}

func (t *SkipGroup) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *SkipGroup) AsJSON() any                { return t.AsJSONOf(t) }
func (t *SkipGroup) AsJSONStr() string          { return t.AsJSONStrOf(t) }
