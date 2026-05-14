package peg

import (
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Sequence struct {
	ModelBase
	Sequence []Model
}

func (s *Sequence) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	var items []Tree
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
	var tree Tree = NIL
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

func (t *Sequence) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Sequence) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Sequence) AsJSONStr() string   { return t.AsJSONStrOf(t) }
