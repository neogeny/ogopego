package peg

import (
	"github.com/neogeny/ogopego/context"
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Sequence struct {
	ModelBase
	Sequence []Model
}

func (s *Sequence) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	var items []trees.Tree
	cutSeen := false
	for _, el := range s.Sequence {
		if _, ok := el.(*Cut); ok {
			cutSeen = true
			ctx.Tracer().TraceCut(ctx)
			continue
		}
		result, err := el.Parse(ctx)
		if err != nil {
			err = context.MarkCut(err, cutSeen)
			ctx.Reset(mark)
			return nil, err
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	var tree trees.Tree = trees.NIL
	switch len(items) {
	case 0:
	case 1:
		tree = items[0]
	default:
		tree = &trees.Seq{Items: items}
	}
	tree.OrCutSeen(cutSeen)
	return tree, nil
}

func (t *Sequence) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Sequence) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Sequence) AsJSONStr() string          { return t.AsJSONStrOf(t) }
