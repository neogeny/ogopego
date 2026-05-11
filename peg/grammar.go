package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Grammar struct {
	ModelBase
	Name       string
	Directives map[string]any
	Keywords   []string
	Rules      []*Rule
	Analyzed   bool
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

func (g *Grammar) ParseTreeFrom(ctx context.Ctx, start string) (trees.Tree, error) {
	if len(g.Keywords) > 0 {
		ctx.SetKeywords(g.Keywords)
	}
	rule, err := g.GetRule(start)
	if err != nil {
		return nil, err
	}
	return rule.Parse(ctx)
}

func (g *Grammar) ParseTree(ctx context.Ctx) (trees.Tree, error) {
	start := "start"
	if _, err := g.GetRule(start); err != nil {
		if len(g.Rules) == 0 {
			return nil, fmt.Errorf("no rules in grammar")
		}
		start = g.Rules[0].Name
	}
	return g.ParseTreeFrom(ctx, start)
}
