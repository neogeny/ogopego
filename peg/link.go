// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// Link resolves rule references within a Box expression.
func (b *Box) Link(rules map[string]*Rule) error {
	if b.Exp != nil {
		return b.Exp.Link(rules)
	}
	return nil
}

// Link resolves rule references within a Join expression.
func (j *Join) Link(rules map[string]*Rule) error {
	if err := j.Exp.Link(rules); err != nil {
		return err
	}
	return j.Sep.Link(rules)
}

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
	c.Target = rule
	return nil
}

// Link resolves rule references within a RuleInclude expression.
func (r *RuleInclude) Link(rules map[string]*Rule) error {
	rule, ok := rules[r.Name]
	if !ok {
		return fmt.Errorf("rule include references undefined rule: %s", r.Name)
	}
	r.Exp = rule.Exp
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
