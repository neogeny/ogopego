// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/neogeny/ogopego/pkg/trees"
)

// comp is a helper struct for compiling grammar trees.
type comp struct {
	path []string
}

// push adds a label to the compilation path.
func (c *comp) push(label string) *comp {
	p := make([]string, len(c.path)+1)
	copy(p, c.path)
	p[len(c.path)] = label
	return &comp{path: p}
}

// error creates a formatted error message with the compilation path.
func (c *comp) error(msg string) error {
	if len(c.path) == 0 {
		return fmt.Errorf("compile: %s", msg)
	}
	return fmt.Errorf("compile: %s at %s", msg, strings.Join(c.path, " -> "))
}

// CompileGrammar compiles a parse tree into a Grammar object.
func CompileGrammar(tree trees.Tree) (*Grammar, error) {
	c := &comp{}
	g, err := c.compileGrammar(tree)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// node extracts the type name and inner tree from a RuleNode.
func (c *comp) node(tree trees.Tree) (string, trees.Tree, error) {
	rn, ok := tree.(*trees.Node)
	if !ok {
		return "", nil, c.error(fmt.Sprintf("expected RuleNode, got %T", tree))
	}
	return rn.TypeName, rn.Tree, nil
}

// nodeCheck checks if a tree is a RuleNode of a specific type.
func (c *comp) nodeCheck(tree trees.Tree, typename string) (trees.Tree, error) {
	name, inner, err := c.node(tree)
	if err != nil {
		return nil, err
	}
	if name != typename {
		return nil, c.error(fmt.Sprintf("expected %s node, got %s", typename, name))
	}
	return inner, nil
}

// mapGet retrieves a value from a MapNode by key.
func (c *comp) mapGet(tree trees.Tree, key string) (trees.Tree, error) {
	mn, ok := tree.(*trees.MapNode)
	if !ok {
		return nil, c.error(fmt.Sprintf("expected MapNode for key %q, got %T", key, tree))
	}
	val, ok := mn.Entries[key]
	if !ok {
		return nil, c.error(fmt.Sprintf("missing key %q", key))
	}
	return val, nil
}

// mapGetDefault retrieves a value from a MapNode by key, returning a default if not found.
func (c *comp) mapGetDefault(tree trees.Tree, key, def string) string {
	mn, ok := tree.(*trees.MapNode)
	if !ok {
		return def
	}
	val, ok := mn.Entries[key]
	if !ok {
		return def
	}
	return textValue(val)
}

// textValue extracts the string value from a Text tree node.
func textValue(tree trees.Tree) string {
	t, ok := tree.(*trees.Text)
	if ok {
		return t.Value
	}
	return ""
}

// listValue extracts a slice of trees from a Seq or List tree node.
func listValue(tree trees.Tree) []trees.Tree {
	switch t := tree.(type) {
	case *trees.Seq:
		return t.Items
	case *trees.List:
		return t.Items
	default:
		return nil
	}
}

// strListValue extracts a slice of strings from a list of Text tree nodes.
func strListValue(tree trees.Tree) []string {
	items := listValue(tree)
	if items == nil {
		return []string{}
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		s := textValue(item)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// strPairsListValue extracts a slice of strings from a list of Text tree nodes.
func strPairsListValue(tree trees.Tree) map[string]string {
	out := map[string]string{}
	items := listValue(tree)
	if items != nil {
		for _, item := range items {
			pair := strListValue(item)
			if len(pair) == 2 {
				out[pair[0]] = pair[1]
			}
		}
	}
	return out
}

// compileGrammar compiles a "Grammar" tree node into a Grammar object.
func (c *comp) compileGrammar(tree trees.Tree) (*Grammar, error) {
	cc := c.push("Grammar")
	inner, err := cc.nodeCheck(tree, "Grammar")
	if err != nil {
		return nil, err
	}

	mn, ok := inner.(*trees.MapNode)
	if !ok {
		return nil, cc.error(fmt.Sprintf("expected MapNode, got %T", inner))
	}

	name := ""
	if n, ok := mn.Entries["name"]; ok {
		name = textValue(n)
	}

	var rules []*Rule

	if rulesTree, ok := mn.Entries["rules"]; ok {
		ruleList := listValue(rulesTree)
		if ruleList == nil {
			return nil, cc.error("rules is not a list")
		}
		for i, rt := range ruleList {
			rc := cc.push(fmt.Sprintf("rules[%d]", i))
			rule, err := rc.compileRule(rt)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
		}
	}

	var directives [][]string
	if dirsTree, ok := mn.Entries["directives"]; ok {
		dirList := listValue(dirsTree)
		for _, d := range dirList {
			dm, dOk := d.(*trees.MapNode)
			if !dOk {
				continue
			}
			n := textValue(dm.Entries["name"])
			v := textValue(dm.Entries["value"])
			if n != "" {
				directives = append(directives, []string{n, v})
				if n == "grammar" && name == "" {
					name = v
				}
			}
		}
	}

	if name == "" {
		name = "__COMPILED__"
	}

	var keywords []string
	if kwTree, ok := mn.Entries["keywords"]; ok {
		kwOuter := listValue(kwTree)
		for _, innerList := range kwOuter {
			for _, kw := range listValue(innerList) {
				s := textValue(kw)
				if s != "" {
					keywords = append(keywords, s)
				}
				if rn, ok := kw.(*trees.Node); ok && rn.TypeName == "Word" {
					s := textValue(rn.Tree)
					if s != "" {
						keywords = append(keywords, s)
					}
				}
			}
		}
	}

	g := &Grammar{
		Name:       name,
		Directives: directives,
		Keywords:   keywords,
		Rules:      rules,
	}
	if err := g.Initialize(); err != nil {
		return nil, cc.error(fmt.Sprintf("initialization: %v", err))
	}
	return g, nil
}

// compileRule compiles a "Rule" tree node into a Rule object.
func (c *comp) compileRule(tree trees.Tree) (*Rule, error) {
	inner, err := c.nodeCheck(tree, "Rule")
	if err != nil {
		return nil, err
	}

	mn, ok := inner.(*trees.MapNode)
	if !ok {
		return nil, c.error(fmt.Sprintf("expected MapNode for Rule, got %T", inner))
	}

	name := textValue(mn.Entries["name"])
	if name == "" {
		return nil, c.error("rule has no name")
	}

	expTree, err := c.mapGet(inner, "exp")
	if err != nil {
		expTree, err = c.mapGet(mn, "exp")
		if err != nil {
			return nil, c.error("rule has no exp")
		}
	}

	exp, err := c.compileExp(expTree)
	if err != nil {
		return nil, err
	}

	decorators := strListValue(mn.Entries["decorators"])

	params := strListValue(mn.Entries["params"])
	kwparams := strPairsListValue(mn.Entries["kwparams"])

	isName := slices.Contains(decorators, "name") ||
		slices.Contains(decorators, "isname")
	noMemo := slices.Contains(decorators, "nomemo")
	noStak := slices.Contains(decorators, "nostak")

	trimmed := strings.TrimLeft(name, "_")
	isTokn := len(trimmed) > 0 && unicode.IsUpper(rune(trimmed[0])) ||
		slices.Contains(decorators, "token") ||
		slices.Contains(decorators, "tokn")

	r := &Rule{
		Exp:        exp,
		Name:       name,
		Params:     params,
		KWParams:   kwparams,
		Decorators: decorators,
		IsName:     isName,
		NoMemo:     noMemo,
		NoStak:     noStak,
		IsTokn:     isTokn,
	}
	return r, nil
}

// compileExp compiles an expression tree node into a Model object.
func (c *comp) compileExp(tree trees.Tree) (Model, error) {
	typename, inner, err := c.node(tree)
	if err != nil {
		return nil, err
	}

	cc := c.push(typename)

	var exp Model
	switch typename {
	case "bool":
		return c.compileExp(inner)

	case "Alert":
		msgTree, err := c.mapGet(inner, "message")
		if err != nil {
			return nil, err
		}
		msg, err := c.compileExp(msgTree)
		if err != nil {
			return nil, err
		}
		cmsg, ok := msg.(*Constant)
		if ok {
			exp = &Alert{Constant: *cmsg}
		} else {
			exp = &Alert{Constant: Constant{Literal: fmt.Sprintf("%v", msg)}}
		}

	case "BasedRule":
		exp = &NULL{}

	case "Box":
		exp = &NULL{}

	case "Call":
		exp = &Call{Name: textValue(inner)}

	case "Choice":
		items := listValue(inner)
		var opts []*Option
		for _, item := range items {
			e, err := cc.compileExp(item)
			if err != nil {
				return nil, err
			}
			if opt, ok := e.(*Option); ok {
				opts = append(opts, opt)
			} else {
				opts = append(opts, &Option{Exp: e})
			}
		}
		exp = &Choice{Options: opts}

	case "Option":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Option{Exp: e}

	case "Closure":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Closure{Exp: e}

	case "Constant":
		exp = &Constant{Literal: textValue(inner)}

	case "Cut":
		exp = &Cut{}

	case "Dot":
		exp = &Dot{}

	case "EOF", "Eof":
		exp = &EOF{}

	case "EOL", "Eol":
		exp = &EOL{}

	case "EmptyClosure":
		exp = &EmptyClosure{}

	case "Fail":
		exp = &Fail{}

	case "Gather":
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		sepTree, err := cc.mapGet(inner, "sep")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		s, err := cc.compileExp(sepTree)
		if err != nil {
			return nil, err
		}
		exp = &Gather{Exp: e, Sep: s}

	case "Grammar":
		exp = &NULL{}

	case "GrammarSemantics":
		exp = &NULL{}

	case "Group":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Group{Exp: e}

	case "Join":
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		sepTree, err := cc.mapGet(inner, "sep")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		s, err := cc.compileExp(sepTree)
		if err != nil {
			return nil, err
		}
		exp = &Join{Exp: e, Sep: s}

	case "Lookahead":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Lookahead{Exp: e}

	case "Model":
		exp = &NULL{}

	case "ModelContext":
		exp = &NULL{}

	case "NULL":
		exp = &NULL{}

	case "Named":
		name := cc.mapGetDefault(inner, "name", "")
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		exp = &Named{Exp: e, Name: name}

	case "NamedBox":
		exp = &NULL{}

	case "NamedList":
		name := cc.mapGetDefault(inner, "name", "")
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		exp = &NamedList{Exp: e, Name: name}

	case "NegativeLookahead":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &NegativeLookahead{Exp: e}

	case "Optional":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Optional{Exp: e}

	case "Override":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &Override{Exp: e}

	case "OverrideList":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &OverrideList{Exp: e}

	case "Pattern":
		exp = &Pattern{Pattern: textValue(inner)}

	case "Patterns":
		var items []trees.Tree
		if t, err := cc.mapGet(inner, "tree"); err == nil {
			items = listValue(t)
		} else {
			items = listValue(inner)
		}
		var exps []Model
		for _, item := range items {
			e, err := cc.compileExp(item)
			if err != nil {
				return nil, err
			}
			exps = append(exps, e)
		}
		if len(exps) == 1 {
			exp = exps[0]
		} else {
			var opts []*Option
			for _, e := range exps {
				opts = append(opts, &Option{Exp: e})
			}
			exp = &Choice{Options: opts}
		}

	case "PositiveClosure":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &PositiveClosure{Exp: e}

	case "PositiveGather":
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		sepTree, err := cc.mapGet(inner, "sep")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		s, err := cc.compileExp(sepTree)
		if err != nil {
			return nil, err
		}
		exp = &PositiveGather{Exp: e, Sep: s}

	case "PositiveJoin", "RightJoin", "LeftJoin":
		expTree, err := cc.mapGet(inner, "exp")
		if err != nil {
			return nil, err
		}
		sepTree, err := cc.mapGet(inner, "sep")
		if err != nil {
			return nil, err
		}
		e, err := cc.compileExp(expTree)
		if err != nil {
			return nil, err
		}
		s, err := cc.compileExp(sepTree)
		if err != nil {
			return nil, err
		}
		exp = &PositiveJoin{Exp: e, Sep: s}

	case "Rule":
		exp = &NULL{}

	case "RuleInclude":
		exp = &RuleInclude{Name: textValue(inner)}

	case "Sequence":
		var items []trees.Tree
		if t, err := cc.mapGet(inner, "tree"); err == nil {
			items = listValue(t)
		} else {
			items = listValue(inner)
		}
		var exps []Model
		for _, item := range items {
			e, err := cc.compileExp(item)
			if err != nil {
				return nil, err
			}
			exps = append(exps, e)
		}
		if len(exps) == 1 {
			exp = exps[0]
		} else {
			exp = &Sequence{Sequence: exps}
		}

	case "SkipGroup":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &SkipGroup{Exp: e}

	case "SkipTo":
		e, err := cc.compileExp(inner)
		if err != nil {
			return nil, err
		}
		exp = &SkipTo{Exp: e}

	case "Synth":
		exp = &NULL{}

	case "Token":
		exp = &Token{Token: textValue(inner)}

	case "Void":
		exp = &Void{}

	default:
		return nil, cc.error(fmt.Sprintf("unknown expression type %q", typename))
	}

	return exp, nil
}
