package peg

import (
	"testing"

	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/tree"
	"github.com/neogeny/ogopego/util/pyre"
)

func ctxFrom(s string) Ctx {
	c := input.NewStrCursor(s)
	pat, err := pyre.Compile(`(?m)[ \t]+`)
	if err != nil {
		panic(err)
	}
	c.SetPatterns(&input.TokenizingPatterns{Wsp: pat})
	return context.NewBaseCtx(c)
}

func TestParseToken(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Token{Token: "hello"}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt, ok := result.(*tree.Text)
	if !ok || tt.Value != "hello" {
		t.Errorf("expected Text{hello}, got %T %+v", result, result)
	}
}

func TestParseTokenFail(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Token{Token: "wrong"}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParsePattern(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Pattern{Pattern: `\w+`}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt, ok := result.(*tree.Text)
	if !ok || tt.Value != "hello" {
		t.Errorf("expected Text{hello}, got %T %+v", result, result)
	}
}

func TestParseSequence(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Sequence{
		Sequence: []Model{
			&Token{Token: "hello"},
			&Token{Token: "world"},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq, ok := result.(*tree.Seq)
	if !ok {
		t.Fatalf("expected Seq, got %T", result)
	}
	if len(seq.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(seq.Items))
	}
	t1 := seq.Items[0].(*tree.Text)
	t2 := seq.Items[1].(*tree.Text)
	if t1.Value != "hello" || t2.Value != "world" {
		t.Errorf("expected hello, world, got %q, %q", t1.Value, t2.Value)
	}
}

func TestParseSequenceEmpty(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Sequence{Sequence: []Model{}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil for empty sequence, got %T", result)
	}
}

func TestParseChoiceFirst(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Choice{
		Options: []*Option{
			{Box: Box{Exp: &Token{Token: "hello"}}},
			{Box: Box{Exp: &Token{Token: "world"}}},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
}

func TestParseChoiceSecond(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Choice{
		Options: []*Option{
			{Box: Box{Exp: &Token{Token: "hello"}}},
			{Box: Box{Exp: &Token{Token: "world"}}},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "world" {
		t.Errorf("expected 'world', got %q", tt.Value)
	}
}

func TestParseChoiceFail(t *testing.T) {
	ctx := ctxFrom("nope")
	expr := &Choice{
		Options: []*Option{
			{Box: Box{Exp: &Token{Token: "hello"}}},
			{Box: Box{Exp: &Token{Token: "world"}}},
		},
	}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error when no option matches")
	}
}

func TestParseOptionalMatches(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Optional{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
}

func TestParseOptionalNoMatch(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Optional{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil for failed optional, got %T", result)
	}
	if ctx.Mark() != 0 {
		t.Errorf("expected cursor at 0 after failed optional, got %d", ctx.Mark())
	}
}

func TestParseClosureMultiple(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &Closure{Box: Box{Exp: &Token{Token: "a"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq, ok := result.(*tree.Seq)
	if !ok {
		t.Fatalf("expected Seq, got %T", result)
	}
	if len(seq.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(seq.Items))
	}
}

func TestParseClosureZero(t *testing.T) {
	ctx := ctxFrom("bbb")
	expr := &Closure{Box: Box{Exp: &Token{Token: "a"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil for zero closure, got %T", result)
	}
}

func TestParsePositiveClosure(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &PositiveClosure{Closure: Closure{Box: Box{Exp: &Token{Token: "a"}}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq, ok := result.(*tree.Seq)
	if !ok {
		t.Fatalf("expected Seq, got %T", result)
	}
	if len(seq.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(seq.Items))
	}
}

func TestParsePositiveClosureFail(t *testing.T) {
	ctx := ctxFrom("bbb")
	expr := &PositiveClosure{Closure: Closure{Box: Box{Exp: &Token{Token: "a"}}}}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error when positive closure can't match at least once")
	}
}

func TestParseGroup(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Group{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
}

func TestParseLookahead(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Lookahead{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
	if ctx.Mark() != 0 {
		t.Errorf("expected cursor restored to 0 after lookahead, got %d", ctx.Mark())
	}
}

func TestParseLookaheadFail(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Lookahead{Box: Box{Exp: &Token{Token: "hello"}}}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseNegativeLookahead(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &NegativeLookahead{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil, got %T", result)
	}
}

func TestParseNegativeLookaheadFail(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &NegativeLookahead{Box: Box{Exp: &Token{Token: "hello"}}}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error when negative lookahead matches")
	}
}

func TestParseNamed(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Named{NamedBox: NamedBox{Box: Box{Exp: &Token{Token: "hello"}}, Name: "greeting"}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	named, ok := result.(*tree.Named)
	if !ok {
		t.Fatalf("expected Named, got %T", result)
	}
	if named.Name != "greeting" {
		t.Errorf("expected name 'greeting', got %q", named.Name)
	}
	tt := named.Value.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
}

func TestParseOverride(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Override{Box: Box{Exp: &Token{Token: "hello"}}}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Override); !ok {
		t.Fatalf("expected Override, got %T", result)
	}
}

func TestParseRule(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Rule{
		NamedBox: NamedBox{
			Box:  Box{Exp: &Token{Token: "hello"}},
			Name: "test",
		},
		Params: []string{"test"},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rn, ok := result.(*tree.TreeNode)
	if !ok {
		t.Fatalf("expected RuleNode, got %T", result)
	}
	if rn.TypeName != "test" {
		t.Errorf("expected 'test', got %q", rn.TypeName)
	}
	tt := rn.Tree.(*tree.Text)
	if tt.Value != "hello" {
		t.Errorf("expected 'hello', got %q", tt.Value)
	}
}

func TestParseGrammar(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "hello"}},
					Name: "start",
				},
				Params: []string{"start"},
			},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rn, ok := result.(*tree.TreeNode)
	if !ok {
		t.Fatalf("expected RuleNode, got %T", result)
	}
	if rn.TypeName != "start" {
		t.Errorf("expected 'start', got %q", rn.TypeName)
	}
}

func TestParseGrammarMultipleRules(t *testing.T) {
	ctx := ctxFrom("hello universe")
	expr := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "hello"}},
					Name: "first",
				},
				Params: []string{"first"},
			},
			{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "universe"}},
					Name: "second",
				},
				Params: []string{"second"},
			},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq, ok := result.(*tree.Seq)
	if !ok {
		t.Fatalf("expected Seq, got %T", result)
	}
	if len(seq.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(seq.Items))
	}
	rn1 := seq.Items[0].(*tree.TreeNode)
	if rn1.TypeName != "first" {
		t.Errorf("expected 'first', got %q", rn1.TypeName)
	}
}

func TestParseEOF(t *testing.T) {
	ctx := ctxFrom("")
	expr := &EOF{}
	_, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEOFFail(t *testing.T) {
	ctx := ctxFrom("not empty")
	expr := &EOF{}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseDot(t *testing.T) {
	ctx := ctxFrom("x")
	expr := &Dot{}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "x" {
		t.Errorf("expected 'x', got %q", tt.Value)
	}
}

func TestParseVoid(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Void{}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil, got %T", result)
	}
}

func TestParseFail(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Fail{}
	_, err := expr.Parse(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseNull(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &NULL{}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*tree.Nil); !ok {
		t.Errorf("expected Nil, got %T", result)
	}
}

func TestParseConstant(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Constant{Literal: "test"}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tt := result.(*tree.Text)
	if tt.Value != "test" {
		t.Errorf("expected 'test', got %q", tt.Value)
	}
}

func TestParseChoiceResetsCursor(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Choice{
		Options: []*Option{
			{Box: Box{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "wrong"}, &Token{Token: "stuff"}},
			}}},
			{Box: Box{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "hello"}, &Token{Token: "world"}},
			}}},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq := result.(*tree.Seq)
	if len(seq.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(seq.Items))
	}
}

func TestParseClosureIncremental(t *testing.T) {
	ctx := ctxFrom("a b c")
	expr := &Sequence{
		Sequence: []Model{
			&Token{Token: "a"},
			&Closure{Box: Box{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "b"}, &Token{Token: "c"}},
			}}},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seq := result.(*tree.Seq)
	if len(seq.Items) != 2 {
		t.Errorf("expected 2 items (a + closure containing b c), got %d", len(seq.Items))
	}
}

func TestParseKeywordIsKeyword(t *testing.T) {
	ctx := ctxFrom("if")
	ctx.SetKeywords([]string{"if", "else"})
	if !ctx.IsKeyword("if") {
		t.Error("expected 'if' to be keyword")
	}
}

func TestParseFoldIntegration(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Sequence{
		Sequence: []Model{
			&Named{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "hello"}},
					Name: "first",
				},
			},
			&Named{
				NamedBox: NamedBox{
					Box:  Box{Exp: &Token{Token: "world"}},
					Name: "second",
				},
			},
		},
	}
	result, err := expr.Parse(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	folded := tree.Fold(result)
	mn, ok := folded.(*tree.MapNode)
	if !ok {
		t.Fatalf("expected MapNode after Fold, got %T", folded)
	}
	if len(mn.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(mn.Entries))
	}
	if mn.Entries["first"] == nil || mn.Entries["second"] == nil {
		t.Errorf("missing keys: first=%v, second=%v", mn.Entries["first"], mn.Entries["second"])
	}
}
