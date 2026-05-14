package peg

import (
	"github.com/neogeny/ogopego/context"
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Group struct {
	Box
}

func (g *Group) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := g.Exp.Parse(ctx)
	if err != nil {
		if pf, ok := err.(*context.ParseFailure); ok {
			pf.CutSeen = false
		}
		return nil, err
	}
	result.TakeCutSeen()
	return result, nil
}

func (t *Group) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Group) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Group) AsJSONStr() string          { return t.AsJSONStrOf(t) }
