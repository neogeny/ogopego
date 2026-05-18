// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

func validateExp(exp Model) error {
	if exp != nil {
		return exp.ValidateLinked()
	}
	return nil
}

func validateSep(exp, sep Model) error {
	if err := validateExp(exp); err != nil {
		return err
	}
	return validateExp(sep)
}

// ValidateLinked checks that all nested expressions are linked.
func (g *Group) ValidateLinked() error             { return validateExp(g.Exp) }
func (o *Optional) ValidateLinked() error          { return validateExp(o.Exp) }
func (c *Closure) ValidateLinked() error           { return validateExp(c.Exp) }
func (p *PositiveClosure) ValidateLinked() error   { return validateExp(p.Exp) }
func (l *Lookahead) ValidateLinked() error         { return validateExp(l.Exp) }
func (n *NegativeLookahead) ValidateLinked() error { return validateExp(n.Exp) }
func (s *SkipGroup) ValidateLinked() error         { return validateExp(s.Exp) }
func (s *SkipTo) ValidateLinked() error            { return validateExp(s.Exp) }
func (o *Override) ValidateLinked() error          { return validateExp(o.Exp) }
func (o *OverrideList) ValidateLinked() error      { return validateExp(o.Exp) }
func (s *Synth) ValidateLinked() error             { return validateExp(s.Exp) }
func (o *Option) ValidateLinked() error            { return validateExp(o.Exp) }
func (n *Named) ValidateLinked() error             { return validateExp(n.Exp) }
func (n *NamedList) ValidateLinked() error         { return validateExp(n.Exp) }
func (r *Rule) ValidateLinked() error              { return validateExp(r.Exp) }

func (j *Join) ValidateLinked() error           { return validateSep(j.Exp, j.Sep) }
func (p *PositiveJoin) ValidateLinked() error   { return validateSep(p.Exp, p.Sep) }
func (g *Gather) ValidateLinked() error         { return validateSep(g.Exp, g.Sep) }
func (p *PositiveGather) ValidateLinked() error { return validateSep(p.Exp, p.Sep) }

// ValidateLinked checks if all expressions in the Sequence are linked.
func (s *Sequence) ValidateLinked() error {
	for _, el := range s.Sequence {
		if err := el.ValidateLinked(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateLinked checks if all options in the Choice are linked.
func (c *Choice) ValidateLinked() error {
	for _, opt := range c.Options {
		if err := opt.ValidateLinked(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateLinked checks if the Call's target rule is linked.
func (c *Call) ValidateLinked() error {
	if c.rule == nil {
		return fmt.Errorf("call to %q is not linked", c.Name)
	}
	return nil
}

// ValidateLinked checks if the RuleInclude's expression is linked.
func (r *RuleInclude) ValidateLinked() error {
	if r.exp == nil {
		return fmt.Errorf("rule include %q is not linked", r.Name)
	}
	return nil
}

// ValidateLinked checks if all rules in the Grammar are linked.
func (g *Grammar) ValidateLinked() error {
	for _, rule := range g.Rules {
		if err := rule.Exp.ValidateLinked(); err != nil {
			return fmt.Errorf("rule %q: %v", rule.Name, err)
		}
	}
	return nil
}
