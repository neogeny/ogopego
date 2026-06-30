// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestPrettyToken(t *testing.T) {
	m := &Token{Token: "hello"}
	got := m.PrettyPrint()
	want := `"hello"`
	assert.Equal(t, want, got)
}

func TestPrettyPattern(t *testing.T) {
	m := &Pattern{Pattern: `\w+`}
	got := m.PrettyPrint()
	want := `/\w+/`
	assert.Equal(t, want, got)
}

func TestPrettyPatternWithSlash(t *testing.T) {
	m := &Pattern{Pattern: `a/b`}
	got := m.PrettyPrint()
	want := `?"a/b"`
	assert.Equal(t, want, got)
}

func TestPrettyConstant(t *testing.T) {
	m := &Constant{Literal: "hello"}
	got := m.PrettyPrint()
	want := "`hello`"
	assert.Equal(t, want, got)
}

func TestPrettyConstantMultiLine(t *testing.T) {
	m := &Constant{Literal: "line1\nline2\nline3"}
	got := m.PrettyPrint()
	want := "```line1\nline2\nline3```"
	assert.Equal(t, want, got)
}

func TestPrettyAlert(t *testing.T) {
	m := &Alert{Constant: Constant{Literal: "warning"}, Level: 2}
	got := m.PrettyPrint()
	want := "^^`warning`"
	assert.Equal(t, want, got)
}

func TestPrettyCall(t *testing.T) {
	m := &Call{Name: "expr"}
	got := m.PrettyPrint()
	assert.Equal(t, "expr", got)
}

func TestPrettyRuleInclude(t *testing.T) {
	m := &RuleInclude{Name: "base_rule"}
	got := m.PrettyPrint()
	want := ">base_rule"
	assert.Equal(t, want, got)
}

func TestPrettyLeafTerminals(t *testing.T) {
	tests := []struct {
		m    Model
		want string
	}{
		{&Cut{}, "~"},
		{&Dot{}, "."},
		{&EOF{}, "$"},
		{&EOL{}, "$->"},
		{&Fail{}, "!()"},
		{&NULL{}, ""},
		{&Void{}, "()"},
		{&EmptyClosure{}, "{}"},
	}
	for _, tt := range tests {
		got := tt.m.PrettyPrint()
		assert.Equal(t, tt.want, got, "%T.PrettyPrint()", tt.m)
	}
}

func TestPrettyGroup(t *testing.T) {
	inner := &Token{Token: "x"}
	m := &Group{Exp: inner}
	got := m.PrettyPrint()
	want := `("x")`
	assert.Equal(t, want, got)
}

func TestPrettyGroupMultiLine(t *testing.T) {
	longTokens := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "iota", "kappa"}
	items := make([]Model, len(longTokens))
	for i, s := range longTokens {
		items[i] = &Token{Token: s}
	}
	seq := &Sequence{Sequence: items}
	m := &Group{Exp: seq}
	got := m.PrettyPrint()
	assert.True(t, strings.HasPrefix(got, "(\n"), "expected multi-line group, got %q", got)
	assert.True(t, strings.HasSuffix(got, ")"), "expected group to end with ')', got %q", got)
}

func TestPrettyOptional(t *testing.T) {
	m := &Optional{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `["x"]`
	assert.Equal(t, want, got)
}

func TestPrettyClosure(t *testing.T) {
	m := &Closure{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `{"x"}*`
	assert.Equal(t, want, got)
}

func TestPrettyPositiveClosure(t *testing.T) {
	m := &PositiveClosure{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `{"x"}+`
	assert.Equal(t, want, got)
}

func TestPrettyLookahead(t *testing.T) {
	m := &Lookahead{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `&"x"`
	assert.Equal(t, want, got)
}

func TestPrettyNegativeLookahead(t *testing.T) {
	m := &NegativeLookahead{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `!"x"`
	assert.Equal(t, want, got)
}

func TestPrettySkipTo(t *testing.T) {
	m := &SkipTo{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `->"x"`
	assert.Equal(t, want, got)
}

func TestPrettySkipGroup(t *testing.T) {
	m := &SkipGroup{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `(?:` + `"x"` + `)`
	assert.Equal(t, want, got)
}

func TestPrettyOverride(t *testing.T) {
	m := &Override{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `="x"`
	assert.Equal(t, want, got)
}

func TestPrettyOverrideList(t *testing.T) {
	m := &OverrideList{Exp: &Token{Token: "x"}}
	got := m.PrettyPrint()
	want := `+="x"`
	assert.Equal(t, want, got)
}

func TestPrettyNamed(t *testing.T) {
	m := &Named{Exp: &Token{Token: "x"}, Name: "value"}
	got := m.PrettyPrint()
	want := `value="x"`
	assert.Equal(t, want, got)
}

func TestPrettyNamedList(t *testing.T) {
	m := &NamedList{Exp: &Token{Token: "x"}, Name: "values"}
	got := m.PrettyPrint()
	want := `values+="x"`
	assert.Equal(t, want, got)
}

func TestPrettySequence(t *testing.T) {
	m := &Sequence{
		Sequence: []Model{
			&Token{Token: "hello"},
			&Token{Token: "world"},
		},
	}
	got := m.PrettyPrint()
	want := `"hello" "world"`
	assert.Equal(t, want, got)
}

func TestPrettyChoice(t *testing.T) {
	m := &Choice{
		Options: []Model{
			&Option{Exp: &Token{Token: "a"}},
			&Option{Exp: &Token{Token: "b"}},
			&Option{Exp: &Token{Token: "c"}},
		},
	}
	got := m.PrettyPrint()
	want := `"a" | "b" | "c"`
	assert.Equal(t, want, got)
}

func TestPrettyJoin(t *testing.T) {
	m := &Join{
		Exp: &Constant{Literal: "x"},
		Sep: &Token{Token: ","},
	}
	got := m.PrettyPrint()
	want := `","%{` + "`x`" + `}*`
	assert.Equal(t, want, got)
}

func TestPrettyGather(t *testing.T) {
	m := &Gather{
		Exp: &Constant{Literal: "x"},
		Sep: &Token{Token: ","},
	}
	got := m.PrettyPrint()
	want := "\",\".{`x`}*"
	assert.Equal(t, want, got)
}

func TestPrettyPositiveGather(t *testing.T) {
	m := &PositiveGather{
		Exp: &Constant{Literal: "x"},
		Sep: &Token{Token: ","},
	}
	got := m.PrettyPrint()
	want := "\",\".{`x`}+"
	assert.Equal(t, want, got)
}

func TestPrettyRule(t *testing.T) {
	r := &Rule{
		Exp: &Choice{
			Options: []Model{
				&Option{Exp: &Token{Token: "hello"}},
				&Option{Exp: &Token{Token: "world"}},
			},
		},
		Name: "greeting",
	}
	got := r.PrettyPrint()
	assert.True(t, strings.Contains(got, "greeting:"), "expected rule name in output, got %q", got)
	assert.True(t, strings.Contains(got, `"hello"`), "expected 'hello' option in output, got %q", got)
	assert.True(t, strings.Contains(got, `"world"`), "expected 'world' option in output, got %q", got)
}

func TestPrettyRuleWithFlags(t *testing.T) {
	r := &Rule{
		Exp:    &Token{Token: "x"},
		Name:   "test",
		NoStak: true,
		NoMemo: true,
	}
	got := r.PrettyPrint()
	assert.True(t, strings.Contains(got, "@nostak"), "expected @nostak, got %q", got)
	assert.True(t, strings.Contains(got, "@nomemo"), "expected @nomemo, got %q", got)
}

func TestPrettyGrammar(t *testing.T) {
	g := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				Exp:  &Token{Token: "hello"},
				Name: "start",
			},
		},
	}
	got := g.PrettyPrint()
	assert.True(t, strings.Contains(got, "@@grammar :: Test"), "expected grammar directive, got %q", got)
	assert.True(t, strings.Contains(got, "start:"), "expected rule, got %q", got)
}

func TestPrettyGrammarWithKeywords(t *testing.T) {
	g := &Grammar{
		Name:     "Test",
		Keywords: []string{"if", "else", "while", "for", "return", "break", "continue", "let", "in"},
		Rules: []*Rule{
			{
				Exp:  &Token{Token: "x"},
				Name: "start",
			},
		},
	}
	got := g.PrettyPrint()
	assert.True(t, strings.Contains(got, "@@keyword"), "expected keyword directive, got %q", got)
}

func TestPrettyGrammarWithDirectives(t *testing.T) {
	g := &Grammar{
		Name: "Test",
		Directives: [][]string{
			{"whitespace", `\s+`},
			{"comments", `#.*`},
		},
		Rules: []*Rule{
			{
				Exp:  &Token{Token: "x"},
				Name: "start",
			},
		},
	}
	got := g.PrettyPrint()
	assert.True(t, strings.Contains(got, "@@whitespace :: /\\s+/"), "expected whitespace directive, got %q", got)
	assert.True(t, strings.Contains(got, "@@comments :: /#.*/"), "expected comments directive, got %q", got)
}

func TestRailroadsToken(t *testing.T) {
	m := &Token{Token: "hello"}
	got := m.Railroads()
	want := `"hello"`
	assert.Equal(t, want, got)
}

func TestRailroadsCall(t *testing.T) {
	m := &Call{Name: "expr"}
	got := m.Railroads()
	assert.Equal(t, "expr", got)
}

func TestRailroadsDot(t *testing.T) {
	m := &Dot{}
	got := m.Railroads()
	if got != " ∀" && got != " ∀ " {
		t.Errorf("got %q, want %q", got, " ∀")
	}
}

func TestRailroadsEOF(t *testing.T) {
	m := &EOF{}
	got := m.Railroads()
	if !strings.Contains(got, "\x03") {
		t.Errorf("expected ETX in EOF railroad, got %q", got)
	}
}

func TestRailroadsEOL(t *testing.T) {
	m := &EOL{}
	got := m.Railroads()
	if !strings.Contains(got, "\x0a") {
		t.Errorf("expected LF in EOL railroad, got %q", got)
	}
}

func TestRailroadsSequence(t *testing.T) {
	m := &Sequence{
		Sequence: []Model{
			&Token{Token: "a"},
			&Token{Token: "b"},
		},
	}
	got := m.Railroads()
	if !strings.Contains(got, "a") || !strings.Contains(got, "b") {
		t.Errorf("expected both tokens in railroad, got %q", got)
	}
}

func TestRailroadsChoice(t *testing.T) {
	m := &Choice{
		Options: []Model{
			&Option{Exp: &Token{Token: "a"}},
			&Option{Exp: &Token{Token: "b"}},
		},
	}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") && !strings.Contains(got, "├─") {
		t.Errorf("expected branch connectors in choice railroad, got %q", got)
	}
}

func TestRailroadsClosure(t *testing.T) {
	m := &Closure{Exp: &Token{Token: "x"}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬→") {
		t.Errorf("expected loop in closure railroad, got %q", got)
	}
}

func TestRailroadsPositiveClosure(t *testing.T) {
	m := &PositiveClosure{Exp: &Token{Token: "x"}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") {
		t.Errorf("expected loop in positive closure railroad, got %q", got)
	}
}

func TestRailroadsOptional(t *testing.T) {
	m := &Optional{Exp: &Token{Token: "x"}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") && !strings.Contains(got, "├─") {
		t.Errorf("expected branch in optional railroad, got %q", got)
	}
}

func TestRailroadsNamed(t *testing.T) {
	m := &Named{
		Exp:  &Token{Token: "x"},
		Name: "val",
	}
	got := m.Railroads()
	if !strings.Contains(got, "val=(") {
		t.Errorf("expected named wrapper, got %q", got)
	}
}

func TestRailroadsLookahead(t *testing.T) {
	m := &Lookahead{Exp: &Token{Token: "x"}}
	got := m.Railroads()
	if !strings.Contains(got, "&[") {
		t.Errorf("expected lookahead wrapper, got %q", got)
	}
}

func TestRailroadsNegativeLookahead(t *testing.T) {
	m := &NegativeLookahead{Exp: &Token{Token: "x"}}
	got := m.Railroads()
	if !strings.Contains(got, "![") {
		t.Errorf("expected negative lookahead wrapper, got %q", got)
	}
}

func TestRailroadsCheckOutputFormat(t *testing.T) {
	// Simple railroad: just a token
	m := &Token{Token: "x"}
	got := m.Railroads()
	if got == "" {
		t.Fatal("expected non-empty railroad")
	}
}

func TestRailroadsGrammar(t *testing.T) {
	g := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				Exp:  &Token{Token: "x"},
				Name: "start",
			},
		},
	}
	got := g.Railroads()
	if !strings.Contains(got, "start") {
		t.Errorf("expected rule name in grammar railroad, got %q", got)
	}
	if !strings.Contains(got, "●─") || !strings.Contains(got, "─■") {
		t.Errorf("expected rule start/end markers, got %q", got)
	}
}
