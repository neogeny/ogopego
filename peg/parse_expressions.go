package peg

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

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

func (t *Token) Parse(ctx Ctx) (trees.Tree, error) {
	matched, err := ctx.Token(t.Token)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (p *Pattern) Parse(ctx Ctx) (trees.Tree, error) {
	matched, err := ctx.Pattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
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

func (s *SkipGroup) Parse(ctx Ctx) (trees.Tree, error) {
	_, err := s.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return trees.NIL, nil
}

func (l *Lookahead) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	_, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return trees.NIL, nil
}

func (n *NegativeLookahead) Parse(ctx Ctx) (trees.Tree, error) {
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
	return trees.NIL, nil
}

func (n *Named) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

func (n *NamedList) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}

func (o *Override) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

func (o *OverrideList) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

func (d *Dot) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	r, err := ctx.Dot()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &trees.Text{Value: string(r)}, nil
}

func (e *EOF) Parse(ctx Ctx) (trees.Tree, error) {
	if !ctx.Eof() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOF"),
		)
	}
	return trees.NIL, nil
}

func (e *EOL) Parse(ctx Ctx) (trees.Tree, error) {
	if !ctx.MatchEOL() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOL"),
		)
	}
	return &trees.Nil{}, nil
}

func (v *Void) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	err := ctx.Void()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return trees.NIL, nil
}

func (f *Fail) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Fail()
}

func (n *NULL) Parse(ctx Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (c *Constant) Parse(ctx Ctx) (trees.Tree, error) {
	return ctx.Constant(c.Literal)
}

func (a *Alert) Parse(ctx Ctx) (trees.Tree, error) {
	return ctx.Constant(a.Literal)
}

func (e *EmptyClosure) Parse(ctx Ctx) (trees.Tree, error) {
	return &trees.List{Items: nil}, nil
}

func (c *Cut) Parse(ctx Ctx) (trees.Tree, error) {
	ctx.Tracer().TraceCut(ctx)
	t := &trees.Nil{}
	t.OrCutSeen(true)
	return t, nil
}

func (o *Option) Parse(ctx Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RuleInclude) Parse(ctx Ctx) (trees.Tree, error) {
	if r.Exp == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("RuleInclude %q has not been resolved", r.Name))
	}
	return r.Exp.Parse(ctx)
}

func (s *Synth) Parse(ctx Ctx) (trees.Tree, error) {
	return s.Exp.Parse(ctx)
}

func (s *SkipTo) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	for {
		_, err := s.Exp.Parse(ctx)
		if err == nil {
			return trees.NIL, nil
		}
		if err != nil {
			ctx.Reset(mark)
			return nil, err
		}

	}
}

func (c *Comment) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp Comment not yet implemented"))
}

func (e *EOLComment) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp EOLComment not yet implemented"))
}

func (r *Rule) Parse(ctx Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	folded := trees.Fold(result)
	if len(r.Params) == 0 || r.Params[0] == "bool" {
		return folded, nil
	}
	return &trees.Node{TypeName: r.Params[0], Tree: folded}, nil
}
