// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

// Link resolves rule references within expression-bearing types.
func (gr *Group) Link(g *Grammar) error            { return gr.Exp.Link(g) }
func (o *Optional) Link(g *Grammar) error          { return o.Exp.Link(g) }
func (c *Closure) Link(g *Grammar) error           { return c.Exp.Link(g) }
func (p *PositiveClosure) Link(g *Grammar) error   { return p.Exp.Link(g) }
func (l *Lookahead) Link(g *Grammar) error         { return l.Exp.Link(g) }
func (n *NegativeLookahead) Link(g *Grammar) error { return n.Exp.Link(g) }
func (s *SkipGroup) Link(g *Grammar) error         { return s.Exp.Link(g) }
func (s *SkipTo) Link(g *Grammar) error            { return s.Exp.Link(g) }
func (o *Override) Link(g *Grammar) error          { return o.Exp.Link(g) }
func (o *OverrideList) Link(g *Grammar) error      { return o.Exp.Link(g) }
func (s *Synth) Link(g *Grammar) error             { return s.Exp.Link(g) }
func (o *Option) Link(g *Grammar) error            { return o.Exp.Link(g) }

func (n *Named) Link(g *Grammar) error {
	return n.Exp.Link(g)
}
func (n *NamedList) Link(g *Grammar) error { return n.Exp.Link(g) }
func (r *Rule) Link(g *Grammar) error      { return r.Exp.Link(g) }

func (j *Join) Link(g *Grammar) error {
	if err := j.Exp.Link(g); err != nil {
		return err
	}
	return j.Sep.Link(g)
}

func (p *PositiveJoin) Link(g *Grammar) error {
	if err := p.Exp.Link(g); err != nil {
		return err
	}
	return p.Sep.Link(g)
}

func (gr *Gather) Link(g *Grammar) error {
	if err := gr.Exp.Link(g); err != nil {
		return err
	}
	return gr.Sep.Link(g)
}

func (p *PositiveGather) Link(g *Grammar) error {
	if err := p.Exp.Link(g); err != nil {
		return err
	}
	return p.Sep.Link(g)
}

// Link resolves rule references within a Sequence expression.
func (s *Sequence) Link(g *Grammar) error {
	for _, el := range s.Sequence {
		if err := el.Link(g); err != nil {
			return err
		}
	}
	return nil
}

// Link resolves rule references within a Choice expression.
func (c *Choice) Link(g *Grammar) error {
	for _, opt := range c.Options {
		if err := opt.Link(g); err != nil {
			return err
		}
	}
	return nil
}

// Link resolves rule references within a Call expression.
func (c *Call) Link(g *Grammar) error {
	rules := g.RuleMap()
	rule, ok := rules[c.Name]
	if !ok {
		return fmt.Errorf("call to undefined rule: %s", c.Name)
	}
	c.rule = rule
	return nil
}

// Link resolves rule references within a RuleInclude expression.
func (r *RuleInclude) Link(g *Grammar) error {
	rules := g.RuleMap()
	rule, ok := rules[r.Name]
	if !ok {
		return fmt.Errorf("rule include references undefined rule: %s", r.Name)
	}
	r.exp = rule.Exp
	return nil
}

func (g *Grammar) LinkGrammar() error {
	for _, rule := range g.Rules {
		if err := rule.Link(g); err != nil {
			return err
		}
	}
	return nil
}

// Link resolves all rule references within the grammar.
func (g *Grammar) Link(other *Grammar) error {
	panic("should never be called")
}
