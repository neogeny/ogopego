package peg

import (
	"strings"
	"testing"

	"github.com/iancoleman/orderedmap"
	asjson "github.com/neogeny/ogopego/json"
)

func TestPrettyToken(t *testing.T) {
	m := &Token{Token: "hello"}
	got := m.PrettyPrint()
	want := `"hello"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyPattern(t *testing.T) {
	m := &Pattern{Pattern: `\w+`}
	got := m.PrettyPrint()
	want := `/\w+/`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyPatternWithSlash(t *testing.T) {
	m := &Pattern{Pattern: `a/b`}
	got := m.PrettyPrint()
	want := `?"a/b"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyConstant(t *testing.T) {
	m := &Constant{Literal: "hello"}
	got := m.PrettyPrint()
	want := "`hello`"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyConstantMultiLine(t *testing.T) {
	m := &Constant{Literal: "line1\nline2\nline3"}
	got := m.PrettyPrint()
	want := "```line1\nline2\nline3```"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyAlert(t *testing.T) {
	m := &Alert{Constant: Constant{Literal: "warning"}, Level: 2}
	got := m.PrettyPrint()
	want := "^^`warning`"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyCall(t *testing.T) {
	m := &Call{Name: "expr"}
	got := m.PrettyPrint()
	if got != "expr" {
		t.Errorf("got %q, want %q", got, "expr")
	}
}

func TestPrettyRuleInclude(t *testing.T) {
	m := &RuleInclude{Name: "base_rule"}
	got := m.PrettyPrint()
	want := ">base_rule"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
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
		if got != tt.want {
			t.Errorf("%T.PrettyPrint() = %q, want %q", tt.m, got, tt.want)
		}
	}
}

func TestPrettyGroup(t *testing.T) {
	inner := &Token{Token: "x"}
	m := &Group{Box: Box{Exp: inner}}
	got := m.PrettyPrint()
	want := `("x")`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyGroupMultiLine(t *testing.T) {
	longTokens := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "iota", "kappa"}
	items := make([]Model, len(longTokens))
	for i, s := range longTokens {
		items[i] = &Token{Token: s}
	}
	seq := &Sequence{Sequence: items}
	m := &Group{Box: Box{Exp: seq}}
	got := m.PrettyPrint()
	if !strings.HasPrefix(got, "(\n") {
		t.Errorf("expected multi-line group, got %q", got)
	}
	if !strings.HasSuffix(got, ")") {
		t.Errorf("expected group to end with ')', got %q", got)
	}
}

func TestPrettyOptional(t *testing.T) {
	m := &Optional{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `["x"]`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyClosure(t *testing.T) {
	m := &Closure{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `{"x"}*`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyPositiveClosure(t *testing.T) {
	m := &PositiveClosure{Closure: Closure{Box: Box{Exp: &Token{Token: "x"}}}}
	got := m.PrettyPrint()
	want := `{"x"}+`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyLookahead(t *testing.T) {
	m := &Lookahead{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `&"x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyNegativeLookahead(t *testing.T) {
	m := &NegativeLookahead{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `!"x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettySkipTo(t *testing.T) {
	m := &SkipTo{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `->"x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettySkipGroup(t *testing.T) {
	m := &SkipGroup{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `(?:` + `"x"` + `)`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyOverride(t *testing.T) {
	m := &Override{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `="x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyOverrideList(t *testing.T) {
	m := &OverrideList{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.PrettyPrint()
	want := `+="x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyNamed(t *testing.T) {
	m := &Named{NamedBox: NamedBox{Box: Box{Exp: &Token{Token: "x"}}, Name: "value"}}
	got := m.PrettyPrint()
	want := `value="x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyNamedList(t *testing.T) {
	m := &NamedList{Named: Named{NamedBox: NamedBox{Box: Box{Exp: &Token{Token: "x"}}, Name: "values"}}}
	got := m.PrettyPrint()
	want := `values+="x"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
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
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyChoice(t *testing.T) {
	m := &Choice{
		Options: []*Option{
			{Box: Box{Exp: &Token{Token: "a"}}},
			{Box: Box{Exp: &Token{Token: "b"}}},
			{Box: Box{Exp: &Token{Token: "c"}}},
		},
	}
	got := m.PrettyPrint()
	want := `"a" | "b" | "c"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyJoin(t *testing.T) {
	m := &Join{
		Box: Box{Exp: &Constant{Literal: "x"}},
		Sep: &Token{Token: ","},
	}
	got := m.PrettyPrint()
	want := `","%{` + "`x`" + `}*`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrettyGather(t *testing.T) {
	m := &Gather{
		Join: Join{
			Box: Box{Exp: &Constant{Literal: "x"}},
			Sep: &Token{Token: ","},
		},
	}
	got := m.PrettyPrint()
	want := "\",\".{`x`}*"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestPrettyPositiveGather(t *testing.T) {
	m := &PositiveGather{
		Gather: Gather{
			Join: Join{
				Box: Box{Exp: &Constant{Literal: "x"}},
				Sep: &Token{Token: ","},
			},
		},
	}
	got := m.PrettyPrint()
	want := "\",\".{`x`}+"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestPrettyRule(t *testing.T) {
	r := &Rule{
		NamedBox: NamedBox{
			Box: Box{Exp: &Choice{
				Options: []*Option{
					{Box: Box{Exp: &Token{Token: "hello"}}},
					{Box: Box{Exp: &Token{Token: "world"}}},
				},
			}},
			Name: "greeting",
		},
	}
	got := r.PrettyPrint()
	if !strings.Contains(got, "greeting:") {
		t.Errorf("expected rule name in output, got %q", got)
	}
	if !strings.Contains(got, `"hello"`) || !strings.Contains(got, `"world"`) {
		t.Errorf("expected options in output, got %q", got)
	}
}

func TestPrettyRuleWithFlags(t *testing.T) {
	r := &Rule{
		NamedBox: NamedBox{
			Box:  Box{Exp: &Token{Token: "x"}},
			Name: "test",
		},
		NoStak: true,
		NoMemo: true,
	}
	got := r.PrettyPrint()
	if !strings.Contains(got, "@nostak") {
		t.Errorf("expected @nostak, got %q", got)
	}
	if !strings.Contains(got, "@nomemo") {
		t.Errorf("expected @nomemo, got %q", got)
	}
}

func TestPrettyGrammar(t *testing.T) {
	g := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "hello"}},
					Name: "start",
				},
			},
		},
	}
	got := g.PrettyPrint()
	if !strings.Contains(got, "@@grammar :: Test") {
		t.Errorf("expected grammar directive, got %q", got)
	}
	if !strings.Contains(got, "start:") {
		t.Errorf("expected rule, got %q", got)
	}
}

func TestPrettyGrammarWithKeywords(t *testing.T) {
	g := &Grammar{
		Name:     "Test",
		Keywords: []string{"if", "else", "while", "for", "return", "break", "continue", "let", "in"},
		Rules: []*Rule{
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "x"}},
					Name: "start",
				},
			},
		},
	}
	got := g.PrettyPrint()
	if !strings.Contains(got, "@@keyword") {
		t.Errorf("expected keyword directive, got %q", got)
	}
}

func TestPrettyGrammarWithDirectives(t *testing.T) {
	g := &Grammar{
		Name: "Test",
		Directives: func() *asjson.OrderedMap {
			om := orderedmap.New()
			om.Set("whitespace", `\s+`)
			om.Set("comments", `#.*`)
			return om
		}(),
		Rules: []*Rule{
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "x"}},
					Name: "start",
				},
			},
		},
	}
	got := g.PrettyPrint()
	if !strings.Contains(got, "@@whitespace :: /\\s+/") {
		t.Errorf("expected whitespace directive, got %q", got)
	}
	if !strings.Contains(got, "@@comments :: /#.*/") {
		t.Errorf("expected comments directive, got %q", got)
	}
}

func TestRailroadsToken(t *testing.T) {
	m := &Token{Token: "hello"}
	got := m.Railroads()
	want := `"hello"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRailroadsCall(t *testing.T) {
	m := &Call{Name: "expr"}
	got := m.Railroads()
	if got != "expr" {
		t.Errorf("got %q, want %q", got, "expr")
	}
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
		Options: []*Option{
			{Box: Box{Exp: &Token{Token: "a"}}},
			{Box: Box{Exp: &Token{Token: "b"}}},
		},
	}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") && !strings.Contains(got, "├─") {
		t.Errorf("expected branch connectors in choice railroad, got %q", got)
	}
}

func TestRailroadsClosure(t *testing.T) {
	m := &Closure{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬→") {
		t.Errorf("expected loop in closure railroad, got %q", got)
	}
}

func TestRailroadsPositiveClosure(t *testing.T) {
	m := &PositiveClosure{Closure: Closure{Box: Box{Exp: &Token{Token: "x"}}}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") {
		t.Errorf("expected loop in positive closure railroad, got %q", got)
	}
}

func TestRailroadsOptional(t *testing.T) {
	m := &Optional{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.Railroads()
	if !strings.Contains(got, "──┬─") && !strings.Contains(got, "├─") {
		t.Errorf("expected branch in optional railroad, got %q", got)
	}
}

func TestRailroadsNamed(t *testing.T) {
	m := &Named{NamedBox: NamedBox{
		Box:  Box{Exp: &Token{Token: "x"}},
		Name: "val",
	}}
	got := m.Railroads()
	if !strings.Contains(got, "val=(") {
		t.Errorf("expected named wrapper, got %q", got)
	}
}

func TestRailroadsLookahead(t *testing.T) {
	m := &Lookahead{Box: Box{Exp: &Token{Token: "x"}}}
	got := m.Railroads()
	if !strings.Contains(got, "&[") {
		t.Errorf("expected lookahead wrapper, got %q", got)
	}
}

func TestRailroadsNegativeLookahead(t *testing.T) {
	m := &NegativeLookahead{Box: Box{Exp: &Token{Token: "x"}}}
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
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "x"}},
					Name: "start",
				},
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
