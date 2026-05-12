package peg

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/tree"
)

func (s *Sequence) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	var items []tree.Tree
	for _, el := range s.Sequence {
		result, err := el.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			return nil, err
		}
		if _, ok := result.(*tree.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &tree.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &tree.Seq{Items: items}, nil
	}
}

func (t *Token) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	matched, err := ctx.Token(t.Token)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Text{Value: matched}, nil
}

func (p *Pattern) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	matched, err := ctx.Pattern(p.Pattern)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Text{Value: matched}, nil
}

func (g *Group) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := g.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return result, nil
}

func (s *SkipGroup) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := s.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return result, nil
}

func (l *Lookahead) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *NegativeLookahead) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, errors.New("negative lookahead matched")
	}
	return &tree.Nil{}, nil
}

func (n *Named) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Named{Name: n.Name, Value: result}, nil
}

func (n *NamedList) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.NamedAsList{Name: n.Name, Value: result}, nil
}

func (o *Override) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Override{Value: result}, nil
}

func (o *OverrideList) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.OverrideAsList{Value: result}, nil
}

func (d *Dot) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	r, err := ctx.Dot()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Text{Value: string(r)}, nil
}

func (e *EOF) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	if !ctx.Eof() {
		ctx.Reset(mark)
		return nil, &ParseError{Pos: ctx.Mark(), Message: "expected EOF"}
	}
	return &tree.Nil{}, nil
}

func (e *EOL) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	if !ctx.MatchEOL() {
		ctx.Reset(mark)
		return nil, &ParseError{Pos: ctx.Mark(), Message: "expected EOL"}
	}
	return &tree.Nil{}, nil
}

func (v *Void) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	err := ctx.Void()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Nil{}, nil
}

func (f *Fail) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	ctx.Reset(mark)
	return nil, ctx.Fail()
}

func (n *NULL) Parse(ctx Ctx) (tree.Tree, error) {
	return &tree.Nil{}, nil
}

func (c *Constant) Parse(ctx Ctx) (tree.Tree, error) {
	return ctx.Constant(c.Literal)
}

func (a *Alert) Parse(ctx Ctx) (tree.Tree, error) {
	return ctx.Constant(a.Literal)
}

func (e *EmptyClosure) Parse(ctx Ctx) (tree.Tree, error) {
	return &tree.Nil{}, nil
}

func (c *Cut) Parse(ctx Ctx) (tree.Tree, error) {
	return &tree.Nil{}, nil
}

func (o *Option) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return result, nil
}

func (r *RuleInclude) Parse(ctx Ctx) (tree.Tree, error) {
	if r.Exp == nil {
		return nil, fmt.Errorf("RuleInclude %q has not been resolved", r.Name)
	}
	mark := ctx.Mark()
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return &tree.Override{Value: result}, nil
}

func (s *Synth) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := s.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return result, nil
}

func (s *SkipTo) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	ctx.Reset(mark)
	return nil, errors.New("SkipTo not yet implemented")
}

func (c *Comment) Parse(ctx Ctx) (tree.Tree, error) {
	return &tree.Nil{}, nil
}

func (e *EOLComment) Parse(ctx Ctx) (tree.Tree, error) {
	return &tree.Nil{}, nil
}

func (r *Rule) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	folded := tree.Fold(result)
	if len(r.Params) == 0 || r.Params[0] == "bool" {
		return folded, nil
	}
	return &tree.TreeNode{TypeName: r.Params[0], Tree: folded}, nil
}

func (g *Grammar) Parse(ctx Ctx) (tree.Tree, error) {
	mark := ctx.Mark()
	if len(g.Keywords) > 0 {
		ctx.SetKeywords(g.Keywords)
	}
	var items []tree.Tree
	for _, rule := range g.Rules {
		result, err := rule.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			return nil, err
		}
		items = append(items, result)
	}
	switch len(items) {
	case 0:
		return &tree.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &tree.Seq{Items: items}, nil
	}
}
