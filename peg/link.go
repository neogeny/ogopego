// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

func linkExp(exp Model, rules map[string]*Rule) error {
	if exp != nil {
		return exp.Link(rules)
	}
	return nil
}

func linkSep(exp, sep Model, rules map[string]*Rule) error {
	if err := linkExp(exp, rules); err != nil {
		return err
	}
	return linkExp(sep, rules)
}

// Link resolves rule references within expression-bearing types.
func (g *Group) Link(rules map[string]*Rule) error             { return linkExp(g.Exp, rules) }
func (o *Optional) Link(rules map[string]*Rule) error          { return linkExp(o.Exp, rules) }
func (c *Closure) Link(rules map[string]*Rule) error           { return linkExp(c.Exp, rules) }
func (p *PositiveClosure) Link(rules map[string]*Rule) error   { return linkExp(p.Exp, rules) }
func (l *Lookahead) Link(rules map[string]*Rule) error         { return linkExp(l.Exp, rules) }
func (n *NegativeLookahead) Link(rules map[string]*Rule) error { return linkExp(n.Exp, rules) }
func (s *SkipGroup) Link(rules map[string]*Rule) error         { return linkExp(s.Exp, rules) }
func (s *SkipTo) Link(rules map[string]*Rule) error            { return linkExp(s.Exp, rules) }
func (o *Override) Link(rules map[string]*Rule) error          { return linkExp(o.Exp, rules) }
func (o *OverrideList) Link(rules map[string]*Rule) error      { return linkExp(o.Exp, rules) }
func (s *Synth) Link(rules map[string]*Rule) error             { return linkExp(s.Exp, rules) }
func (o *Option) Link(rules map[string]*Rule) error            { return linkExp(o.Exp, rules) }
func (n *Named) Link(rules map[string]*Rule) error             { return linkExp(n.Exp, rules) }
func (n *NamedList) Link(rules map[string]*Rule) error         { return linkExp(n.Exp, rules) }
func (r *Rule) Link(rules map[string]*Rule) error              { return linkExp(r.Exp, rules) }

func (j *Join) Link(rules map[string]*Rule) error           { return linkSep(j.Exp, j.Sep, rules) }
func (p *PositiveJoin) Link(rules map[string]*Rule) error   { return linkSep(p.Exp, p.Sep, rules) }
func (g *Gather) Link(rules map[string]*Rule) error         { return linkSep(g.Exp, g.Sep, rules) }
func (p *PositiveGather) Link(rules map[string]*Rule) error { return linkSep(p.Exp, p.Sep, rules) }

// Link resolves rule references within a Sequence expression.
func (s *Sequence) Link(rules map[string]*Rule) error {
	for _, el := range s.Sequence {
		if err := el.Link(rules); err != nil {
			return err
		}
	}
	return nil
}

// Link resolves rule references within a Choice expression.
func (c *Choice) Link(rules map[string]*Rule) error {
	for _, opt := range c.Options {
		if err := opt.Link(rules); err != nil {
			return err
		}
	}
	return nil
}

// Link resolves rule references within a Call expression.
func (c *Call) Link(rules map[string]*Rule) error {
	rule, ok := rules[c.Name]
	if !ok {
		return fmt.Errorf("call to undefined rule: %s", c.Name)
	}
	c.rule = rule
	return nil
}

// Link resolves rule references within a RuleInclude expression.
func (r *RuleInclude) Link(rules map[string]*Rule) error {
	rule, ok := rules[r.Name]
	if !ok {
		return fmt.Errorf("rule include references undefined rule: %s", r.Name)
	}
	r.exp = rule.Exp
	return nil
}

func (g *Grammar) LinkGrammar() error {
	return g.Link(g.RuleMap())
}

// Link resolves all rule references within the grammar.
func (g *Grammar) Link(rules map[string]*Rule) error {
	for _, rule := range g.Rules {
		if err := rule.Exp.Link(rules); err != nil {
			return err
		}
	}
	return nil
}
