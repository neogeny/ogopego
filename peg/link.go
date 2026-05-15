// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

func (b *Box) Link(rules map[string]*Rule) error {
	if b.Exp != nil {
		return b.Exp.Link(rules)
	}
	return nil
}

func (j *Join) Link(rules map[string]*Rule) error {
	if err := j.Exp.Link(rules); err != nil {
		return err
	}
	return j.Sep.Link(rules)
}

func (s *Sequence) Link(rules map[string]*Rule) error {
	for _, el := range s.Sequence {
		if err := el.Link(rules); err != nil {
			return err
		}
	}
	return nil
}

func (c *Choice) Link(rules map[string]*Rule) error {
	for _, opt := range c.Options {
		if err := opt.Link(rules); err != nil {
			return err
		}
	}
	return nil
}

func (c *Call) Link(rules map[string]*Rule) error {
	rule, ok := rules[c.Name]
	if !ok {
		return fmt.Errorf("call to undefined rule: %s", c.Name)
	}
	c.Target = rule
	return nil
}

func (r *RuleInclude) Link(rules map[string]*Rule) error {
	rule, ok := rules[r.Name]
	if !ok {
		return fmt.Errorf("rule include references undefined rule: %s", r.Name)
	}
	r.Exp = rule.Exp
	return nil
}

func (g *Grammar) Link() error {
	rules := make(map[string]*Rule, len(g.Rules))
	for _, rule := range g.Rules {
		rules[rule.Name] = rule
	}
	for _, rule := range g.Rules {
		if err := rule.Exp.Link(rules); err != nil {
			return err
		}
	}
	return nil
}
