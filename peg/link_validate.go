// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// ValidateLinked checks if the Box's expression is linked.
func (b *Box) ValidateLinked() error {
	if b.Exp != nil {
		return b.Exp.ValidateLinked()
	}
	return nil
}

// ValidateLinked checks if the Join's expressions are linked.
func (j *Join) ValidateLinked() error {
	if err := j.Exp.ValidateLinked(); err != nil {
		return err
	}
	return j.Sep.ValidateLinked()
}

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
	if c.Target == nil {
		return fmt.Errorf("call to %q is not linked", c.Name)
	}
	return nil
}

// ValidateLinked checks if the RuleInclude's expression is linked.
func (r *RuleInclude) ValidateLinked() error {
	if r.Exp == nil {
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
