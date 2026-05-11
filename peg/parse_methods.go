package peg

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

func (s *Sequence) Parse(ctx context.Ctx) (trees.Tree, error) {
	var items []trees.Tree
	for _, el := range s.Sequence {
		result, err := el.Parse(ctx)
		if err != nil {
			return nil, err
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}

func (c *Choice) Parse(ctx context.Ctx) (trees.Tree, error) {
	var lastErr error
	for _, opt := range c.Options {
		mark := ctx.Mark()
		result, err := opt.Parse(ctx)
		if err == nil {
			return result, nil
		}
		ctx.Reset(mark)
		lastErr = err
	}
	if lastErr == nil {
		return nil, errors.New("no option matched")
	}
	return nil, lastErr
}

func (o *Option) Parse(ctx context.Ctx) (trees.Tree, error) {
	return o.Exp.Parse(ctx)
}

func (c *Cut) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (t *Token) Parse(ctx context.Ctx) (trees.Tree, error) {
	matched, err := ctx.Token(t.Token)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (p *Pattern) Parse(ctx context.Ctx) (trees.Tree, error) {
	matched, err := ctx.Pattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (g *Group) Parse(ctx context.Ctx) (trees.Tree, error) {
	return g.Exp.Parse(ctx)
}

func (s *SkipGroup) Parse(ctx context.Ctx) (trees.Tree, error) {
	return s.Exp.Parse(ctx)
}

func (o *Optional) Parse(ctx context.Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return &trees.Nil{}, nil
	}
	return result, nil
}

func (c *Closure) Parse(ctx context.Ctx) (trees.Tree, error) {
	var items []trees.Tree
	for {
		mark := ctx.Mark()
		result, err := c.Exp.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			break
		}
		if ctx.Mark() == mark {
			break
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}

func (p *PositiveClosure) Parse(ctx context.Ctx) (trees.Tree, error) {
	first, err := p.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	var items []trees.Tree
	if _, ok := first.(*trees.Nil); !ok {
		items = append(items, first)
	}
	for {
		mark := ctx.Mark()
		result, err := p.Exp.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			break
		}
		if ctx.Mark() == mark {
			break
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}

func (l *Lookahead) Parse(ctx context.Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	result, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *NegativeLookahead) Parse(ctx context.Ctx) (trees.Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, errors.New("negative lookahead matched")
	}
	return &trees.Nil{}, nil
}

func (n *Named) Parse(ctx context.Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

func (n *NamedList) Parse(ctx context.Ctx) (trees.Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}

func (o *Override) Parse(ctx context.Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

func (o *OverrideList) Parse(ctx context.Ctx) (trees.Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

func (d *Dot) Parse(ctx context.Ctx) (trees.Tree, error) {
	r, err := ctx.Dot()
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: string(r)}, nil
}

func (e *EOF) Parse(ctx context.Ctx) (trees.Tree, error) {
	if !ctx.Eof() {
		return nil, &context.ParseError{Pos: ctx.Mark(), Message: "expected EOF"}
	}
	return &trees.Nil{}, nil
}

func (e *EOL) Parse(ctx context.Ctx) (trees.Tree, error) {
	if !ctx.MatchEOL() {
		return nil, &context.ParseError{Pos: ctx.Mark(), Message: "expected EOL"}
	}
	return &trees.Nil{}, nil
}

func (v *Void) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, ctx.Void()
}

func (f *Fail) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, ctx.Fail()
}

func (n *NULL) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (c *Constant) Parse(ctx context.Ctx) (trees.Tree, error) {
	return ctx.Constant(c.Literal)
}

func (a *Alert) Parse(ctx context.Ctx) (trees.Tree, error) {
	return ctx.Constant(a.Literal)
}

func (e *EmptyClosure) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (r *RuleInclude) Parse(ctx context.Ctx) (trees.Tree, error) {
	if r.Exp == nil {
		return nil, fmt.Errorf("RuleInclude %q has not been resolved", r.Name)
	}
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

func (s *Synth) Parse(ctx context.Ctx) (trees.Tree, error) {
	return s.Exp.Parse(ctx)
}

func (c *Call) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, fmt.Errorf("rule calls not yet supported: %s", c.Name)
}

func (s *SkipTo) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, errors.New("SkipTo not yet implemented")
}

func (j *Join) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, errors.New("Join not yet implemented")
}

func (p *PositiveJoin) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, errors.New("PositiveJoin not yet implemented")
}

func (g *Gather) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, errors.New("Gather not yet implemented")
}

func (p *PositiveGather) Parse(ctx context.Ctx) (trees.Tree, error) {
	return nil, errors.New("PositiveGather not yet implemented")
}

func (c *Comment) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (e *EOLComment) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (p *Patterns) Parse(ctx context.Ctx) (trees.Tree, error) {
	return &trees.Nil{}, nil
}

func (r *Rule) Parse(ctx context.Ctx) (trees.Tree, error) {
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.RuleNode{TypeName: r.Name, Tree: result}, nil
}

func (g *Grammar) Parse(ctx context.Ctx) (trees.Tree, error) {
	if len(g.Keywords) > 0 {
		ctx.SetKeywords(g.Keywords)
	}
	var items []trees.Tree
	for _, rule := range g.Rules {
		result, err := rule.Parse(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, result)
	}
	switch len(items) {
	case 0:
		return &trees.Nil{}, nil
	case 1:
		return items[0], nil
	default:
		return &trees.Seq{Items: items}, nil
	}
}
