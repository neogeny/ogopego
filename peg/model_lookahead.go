package peg

import "fmt"

type Lookahead struct {
	Box
}

type NegativeLookahead struct {
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

func (n *NegativeLookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, ctx.Failure(
			mark,
			fmt.Errorf(
				"negative lookahead matched:%v",
				n.Exp,
			),
		)
	}
	return NIL, nil
}

func (t *Lookahead) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Lookahead) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Lookahead) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (t *NegativeLookahead) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NegativeLookahead) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NegativeLookahead) AsJSONStr() string   { return t.AsJSONStrOf(t) }
