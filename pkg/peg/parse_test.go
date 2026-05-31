// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/trees"
	"github.com/neogeny/ogopego/pkg/util/pyre"
)

func ctxFrom(s string) Ctx {
	c := input.NewStrCursor(s)
	pat, err := pyre.Compile(`(?m)[ \t]+`)
	if err != nil {
		panic(err)
	}
	c.SetPatterns(&input.TokenizingPatterns{Wsp: pat})
	return context.NewCtx(c, nil)
}

func TestParseToken(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Token{Token: "hello"}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt, ok := result.(*trees.Text)
	assert.True(t, ok, "expected Text{hello}, got %T %+v", result, result)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseTokenFail(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Token{Token: "wrong"}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error")
}

func TestParsePattern(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Pattern{Pattern: `\w+`}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt, ok := result.(*trees.Text)
	assert.True(t, ok, "expected Text{hello}, got %T %+v", result, result)
	assert.Equal(t, "hello", tt.Value)
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
	assert.NoError(t, err)
	seq, ok := result.(*trees.Seq)
	assert.True(t, ok, "expected Seq, got %T", result)
	assert.Equal(t, 2, len(seq.Items), "expected 2 items, got %d", len(seq.Items))
	t1 := seq.Items[0].(*trees.Text)
	t2 := seq.Items[1].(*trees.Text)
	assert.Equal(t, "hello", t1.Value)
	assert.Equal(t, "world", t2.Value)
}

func TestParseSequenceEmpty(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Sequence{Sequence: []Model{}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil for empty sequence, got %T", result)
}

func TestParseChoiceFirst(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Choice{
		Options: []*Option{
			{Exp: &Token{Token: "hello"}},
			{Exp: &Token{Token: "world"}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseChoiceSecond(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Choice{
		Options: []*Option{
			{Exp: &Token{Token: "hello"}},
			{Exp: &Token{Token: "world"}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "world", tt.Value)
}

func TestParseChoiceFail(t *testing.T) {
	ctx := ctxFrom("nope")
	expr := &Choice{
		Options: []*Option{
			{Exp: &Token{Token: "hello"}},
			{Exp: &Token{Token: "world"}},
		},
	}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error when no option matches")
}

func TestParseOptionalMatches(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Optional{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseOptionalNoMatch(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Optional{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil for failed optional, got %T", result)
	assert.Equal(t, 0, ctx.Mark(), "expected cursor at 0 after failed optional, got %d", ctx.Mark())
}

func TestParseClosureMultiple(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &Closure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.(*trees.Array)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 3, len(lst.Items))
}

func TestParseClosureZero(t *testing.T) {
	ctx := ctxFrom("bbb")
	expr := &Closure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.(*trees.Array)
	assert.True(t, ok, "expected List for zero closure, got %T", result)
	assert.Equal(t, 0, len(lst.Items))
}

func TestParsePositiveClosure(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &PositiveClosure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.(*trees.Array)
	assert.True(t, ok, "expected List, got %T", result)
	assert.Equal(t, 3, len(lst.Items))
}

func TestParsePositiveClosureFail(t *testing.T) {
	ctx := ctxFrom("bbb")
	expr := &PositiveClosure{Exp: &Token{Token: "a"}}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error when positive closure can't match at least once")
}

func TestParseGroup(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Group{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseLookahead(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Lookahead{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil after lookahead, got %T", result)
	assert.Equal(t, 0, ctx.Mark(), "expected cursor restored to 0 after lookahead, got %d", ctx.Mark())
}

func TestParseLookaheadFail(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Lookahead{Exp: &Token{Token: "hello"}}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error")
}

func TestParseNegativeLookahead(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &NegativeLookahead{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil, got %T", result)
}

func TestParseNegativeLookaheadFail(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &NegativeLookahead{Exp: &Token{Token: "hello"}}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error when negative lookahead matches")
}

func TestParseNamed(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Named{Exp: &Token{Token: "hello"}, Name: "greeting"}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	named, ok := result.(*trees.Named)
	assert.True(t, ok, "expected Named, got %T", result)
	assert.Equal(t, "greeting", named.Name)
	tt := named.Value.(*trees.Text)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseOverride(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Override{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Override)
	assert.True(t, ok, "expected Override, got %T", result)
}

func TestParseRule(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Rule{
		Exp:    &Token{Token: "hello"},
		Name:   "test",
		Params: []string{"test"},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	rn, ok := result.(*trees.Node)
	assert.True(t, ok, "expected RuleNode, got %T", result)
	assert.Equal(t, "test", rn.TypeName)
	tt := rn.Tree.(*trees.Text)
	assert.Equal(t, "hello", tt.Value)
}

func TestParseGrammar(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				Exp:    &Token{Token: "hello"},
				Name:   "start",
				Params: []string{"start"},
			},
		},
	}
	result, err := expr.ParseAt(ctx, nil)
	assert.NoError(t, err)
	rn, ok := result.(*trees.Node)
	assert.True(t, ok, "expected RuleNode, got %T", result)
	assert.Equal(t, "start", rn.TypeName)
}

func TestParseGrammarMultipleRules(t *testing.T) {
	ctx := ctxFrom("hello universe")
	expr := &Grammar{
		Name: "Test",
		Rules: []*Rule{
			{
				Exp:    &Token{Token: "hello"},
				Name:   "first",
				Params: []string{"first"},
			},
			{
				Exp:    &Token{Token: "universe"},
				Name:   "second",
				Params: []string{"second"},
			},
		},
	}
	result, err := expr.ParseAt(ctx, nil)
	assert.NoError(t, err)
	rn, ok := result.(*trees.Node)
	assert.True(t, ok, "expected Node, got %T", result)
	assert.Equal(t, "first", rn.TypeName)
}

func TestParseEOF(t *testing.T) {
	ctx := ctxFrom("")
	expr := &EOF{}
	_, err := expr.Parse(ctx)
	assert.NoError(t, err)
}

func TestParseEOFFail(t *testing.T) {
	ctx := ctxFrom("not empty")
	expr := &EOF{}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error")
}

func TestParseDot(t *testing.T) {
	ctx := ctxFrom("x")
	expr := &Dot{}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "x", tt.Value)
}

func TestParseVoid(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Void{}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil, got %T", result)
}

func TestParseFail(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Fail{}
	_, err := expr.Parse(ctx)
	assert.Error(t, err, "expected error")
}

func TestParseNull(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &NULL{}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	_, ok := result.(*trees.Nil)
	assert.True(t, ok, "expected Nil, got %T", result)
}

func TestParseConstant(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Constant{Literal: "test"}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	tt := result.(*trees.Text)
	assert.Equal(t, "test", tt.Value)
}

func TestParseChoiceResetsCursor(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Choice{
		Options: []*Option{
			{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "wrong"}, &Token{Token: "stuff"}},
			}},
			{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "hello"}, &Token{Token: "world"}},
			}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	seq := result.(*trees.Seq)
	assert.Equal(t, 2, len(seq.Items))
}

func TestParseClosureIncremental(t *testing.T) {
	ctx := ctxFrom("a b c")
	expr := &Sequence{
		Sequence: []Model{
			&Token{Token: "a"},
			&Closure{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "b"}, &Token{Token: "c"}},
			}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	seq := result.(*trees.Seq)
	assert.Equal(t, 2, len(seq.Items), "expected 2 items (a + closure containing b c), got %d", len(seq.Items))
}

func TestParseKeywordIsKeyword(t *testing.T) {
	ctx := ctxFrom("if")
	ctx.Configure(input.Cfg{Keywords: []string{"if", "else"}})
	assert.True(t, ctx.IsKeyword("if"), "expected 'if' to be keyword")
}

func TestParseFoldIntegration(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Sequence{
		Sequence: []Model{
			&Named{
				Exp:  &Token{Token: "hello"},
				Name: "first",
			},
			&Named{
				Exp:  &Token{Token: "world"},
				Name: "second",
			},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	folded := trees.Fold(result)
	mn, ok := folded.(*trees.MapNode)
	assert.True(t, ok, "expected MapNode after Fold, got %T", folded)
	assert.Equal(t, 2, len(mn.Entries))
	assert.NotZero(t, mn.Entries["first"], "missing key 'first'")
	assert.NotZero(t, mn.Entries["second"], "missing key 'second'")
}
