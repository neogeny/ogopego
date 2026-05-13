package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

type Grammar struct {
	ModelBase
	Name       string
	cfg        Cfg
	Directives map[string]any
	Keywords   []string
	Rules      []*Rule
	Analyzed   bool
}

func (g *Grammar) CfgFromDirectives() *Cfg {
	c := Cfg{}

	for name, val := range g.Directives {
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
				c.Whitespace = nil
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

func (g *Grammar) Initialize() error {
	if err := g.Link(); err != nil {
		return err
	}
	if err := g.ValidateLinked(); err != nil {
		return err
	}
	g.Analyzed = true
	return nil
}

func (g *Grammar) GetRule(name string) (*Rule, error) {
	for _, rule := range g.Rules {
		if rule.Name == name {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("rule %q not found", name)
}

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
	return rule.Parse(ctx)
}
