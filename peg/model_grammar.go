// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/trees"
)

// Grammar represents a parsed PEG grammar, containing rules, directives, and keywords.
type Grammar struct {
	ModelBase
	Name       string
	Directives [][]string
	Keywords   []string
	Rules      []*Rule
	Analyzed   bool
	// Allow setting semantics filtering at the grammar definition level
	// Can be used on the boot grammar to implement grammar-parser semantics
	Semantics config.SemanticsFunc
}

// CfgFromDirectives creates a Cfg object based on the grammar's directives.
func (g *Grammar) CfgFromDirectives() *Cfg {
	c := Cfg{}

	for _, d := range g.Directives {
		name := d[0]
		s := d[1]

		switch name {
		case "name":
			c.Name = s
		case "source":
			c.Source = s
		case "start":
			c.Start = s
		case "grammar":
			c.Grammar = s
		case "whitespace":
			if s == "" || s == "None" || s == "False" {
				c.Whitespace = new("")
			} else {
				c.Whitespace = &s
			}
		case "comments":
			c.Comments = s
		case "eol_comments":
			c.EolComments = s
		case "ignorecase":
			c.IgnoreCase = s == "True" || s == "true" || s == "1"
		case "namechars":
			c.NameChars = s
		case "nameguard":
			c.NameGuard = s == "True" || s == "true" || s == "1"
		case "parseinfo":
			c.ParseInfo = s == "True" || s == "true" || s == "1"
		case "trace":
			c.Trace = s == "True" || s == "true" || s == "1"
		case "left_recursion":
			c.NoLeftRecursion = s != "True" && s != "true" && s != "1"
		case "nomemo":
			c.NoMemo = s == "True" || s == "true" || s == "1"
		case "noprunememosoncut":
			c.NoPruneMemosOnCut = s == "True" || s == "true" || s == "1"
		}
	}
	if g.Semantics != nil {
		c.Semantics = g.Semantics
	}
	return &c
}

func (g *Grammar) RuleMap() map[string]*Rule {
	m := make(map[string]*Rule, len(g.Rules))
	for _, rule := range g.Rules {
		m[rule.Name] = rule
	}
	return m
}

// Initialize links and validates the grammar, and marks left-recursive rules.
func (g *Grammar) Initialize() error {
	if err := g.LinkGrammar(); err != nil {
		return err
	}
	if err := g.ValidateLinked(); err != nil {
		return err
	}
	g.markLeftRecursion()
	g.computeAnalysis()
	// NOTE
	//  Do not Optimize here so comparisons with sibling output is possible
	//  Optimization can be postoned until model is used to Parse()
	//g.Optimize()
	g.Analyzed = true
	return nil
}

func (g *Grammar) computeAnalysis() {
	for _, rule := range g.Rules {
		computeLA(rule.Exp)
	}
}

// GetRule retrieves a rule by its name.
func (g *Grammar) GetRule(name string) (*Rule, error) {
	for _, rule := range g.Rules {
		if rule.Name == name {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("rule %q not found", name)
}

// ParseAt parses the input using the grammar, starting from the specified rule.
func (g *Grammar) ParseAt(ctx Ctx, cfg *Cfg) (trees.Tree, error) {
	acfg := g.CfgFromDirectives().New()
	acfg = acfg.Override(cfg)
	acfg.Keywords = g.Keywords
	ctx.Configure(acfg)

	start := acfg.Start
	if start == "" {
		start = "start"
	}
	rule, err := g.GetRule(start)
	if err != nil {
		if len(g.Rules) == 0 {
			return nil, ctx.Failure(
				ctx.Mark(),
				fmt.Errorf("no rules in grammar"),
			)
		}
		rule = g.Rules[0]
	}
	result, err := rule.Parse(ctx)
	if err != nil {
		if dis := ctx.FurthestFailure(); dis != nil {
			return nil, dis
		}
		return nil, err
	}
	return result, nil
}

// normalize Bring fields to a consistent state, setting defaults
// for missing values and ensuring internal consistency.
func (g *Grammar) normalize() {
	for _, r := range g.Rules {
		r.normalize()
	}
}
