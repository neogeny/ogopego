package peg

import "fmt"

func (b *Box) ValidateLinked() error {
	if b.Exp != nil {
		return b.Exp.ValidateLinked()
	}
	return nil
}

func (j *Join) ValidateLinked() error {
	if err := j.Exp.ValidateLinked(); err != nil {
		return err
	}
	return j.Sep.ValidateLinked()
}

func (s *Sequence) ValidateLinked() error {
	for _, el := range s.Sequence {
		if err := el.ValidateLinked(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Choice) ValidateLinked() error {
	for _, opt := range c.Options {
		if err := opt.ValidateLinked(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Call) ValidateLinked() error {
	if c.Target == nil {
		return fmt.Errorf("call to %q is not linked", c.Name)
	}
	return nil
}

func (r *RuleInclude) ValidateLinked() error {
	if r.Exp == nil {
		return fmt.Errorf("rule include %q is not linked", r.Name)
	}
	return nil
}

func (g *Grammar) ValidateLinked() error {
	for _, rule := range g.Rules {
		if err := rule.Exp.ValidateLinked(); err != nil {
			return fmt.Errorf("rule %q: %v", rule.Name, err)
		}
	}
	return nil
}
