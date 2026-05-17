// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"encoding/json"
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// Grammar represents a parsed PEG grammar, containing rules, directives, and keywords.
type Grammar struct {
	ModelBase
	Name       string
	Directives *asjson.OrderedMap
	Keywords   []string
	Rules      []*Rule
	Analyzed   bool
}

// CfgFromDirectives creates a Cfg object based on the grammar's directives.
func (g *Grammar) CfgFromDirectives() *Cfg {
	c := Cfg{}

	if g.Directives == nil {
		return &c
	}

	for _, name := range g.Directives.Keys() {
		val, _ := g.Directives.Get(name)
		s, ok := val.(string)
		if !ok {
			continue
		}

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
	return &c
}

// Initialize links and validates the grammar, and marks left-recursive rules.
func (g *Grammar) Initialize() error {
	if err := g.Link(); err != nil {
		return err
	}
	if err := g.ValidateLinked(); err != nil {
		return err
	}
	g.markLeftRecursion()
	g.Analyzed = true
	return nil
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

// Parse parses the input using the grammar, starting from the specified rule.
func (g *Grammar) Parse(ctx Ctx, cfg *Cfg) (trees.Tree, error) {
	acfg := g.CfgFromDirectives().New()
	acfg = acfg.Override(cfg)
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

// PubMap returns an ordered map of the Grammar's public fields.
func (g *Grammar) PubMap() *asjson.OrderedMap { return util.PubMapOf(g) }

// AsJSON returns a JSON-compatible representation of the Grammar.
func (g *Grammar) AsJSON() any { return g.AsJSONOf(g) }

// AsJSONStr returns a JSON string representation of the Grammar.
func (g *Grammar) AsJSONStr() string { return g.AsJSONStrOf(g) }

// MarshalJSON marshals the Grammar to JSON.
func (g *Grammar) MarshalJSON() ([]byte, error) { return json.Marshal(g.AsJSON()) }
