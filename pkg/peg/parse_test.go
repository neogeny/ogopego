// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/asjson"
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
	assert.Equal(t, "hello", asjson.AsJSON(result))
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
	assert.Equal(t, "hello", asjson.AsJSON(result))
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
	lst := asjson.AsJSON(result).([]any)
	assert.Equal(t, 2, len(lst), "expected 2 items, got %d", len(lst))
	assert.Equal(t, "hello", lst[0].(string))
	assert.Equal(t, "world", lst[1].(string))
}

func TestParseSequenceEmpty(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Sequence{Sequence: []Model{}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.True(t, result == nil, "expected nil for empty sequence, got %T", result)
}

func TestParseChoiceFirst(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Choice{
		Options: []Model{
			&Option{Exp: &Token{Token: "hello"}},
			&Option{Exp: &Token{Token: "world"}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hello", asjson.AsJSON(result))
}

func TestParseChoiceSecond(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Choice{
		Options: []Model{
			&Option{Exp: &Token{Token: "hello"}},
			&Option{Exp: &Token{Token: "world"}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "world", asjson.AsJSON(result))
}

func TestParseChoiceFail(t *testing.T) {
	ctx := ctxFrom("nope")
	expr := &Choice{
		Options: []Model{
			&Option{Exp: &Token{Token: "hello"}},
			&Option{Exp: &Token{Token: "world"}},
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
	assert.Equal(t, "hello", asjson.AsJSON(result))
}

func TestParseOptionalNoMatch(t *testing.T) {
	ctx := ctxFrom("world")
	expr := &Optional{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.True(t, result == nil, "expected nil for failed optional, got %T", result)
	assert.Equal(t, 0, ctx.Mark(), "expected cursor at 0 after failed optional, got %d", ctx.Mark())
}

func TestParseClosureMultiple(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &Closure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.([]any)
	assert.True(t, ok, "expected list, got %T", result)
	assert.Equal(t, 3, len(lst))
}

func TestParseClosureZero(t *testing.T) {
	ctx := ctxFrom("bbb")
	expr := &Closure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.([]any)
	assert.True(t, ok, "expected list for zero closure, got %T", result)
	assert.Equal(t, 0, len(lst))
}

func TestParsePositiveClosure(t *testing.T) {
	ctx := ctxFrom("aaa")
	expr := &PositiveClosure{Exp: &Token{Token: "a"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst, ok := result.([]any)
	assert.True(t, ok, "expected list, got %T", result)
	assert.Equal(t, 3, len(lst))
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
	assert.Equal(t, "hello", asjson.AsJSON(result))
}

func TestParseLookahead(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Lookahead{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.True(t, result == nil, "expected nil after lookahead, got %T", result)
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
	assert.True(t, result == nil, "expected nil, got %T", result)
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
	f := trees.Fold(result)
	m := asjson.AsJSON(f).(map[string]any)
	assert.NotEqual(t, nil, m, "expected non-nil result %v", m)
	assert.Equal(t, "hello", m["greeting"], "expected [:greeting] to be 'hello' %v", m)
}

func TestParseOverride(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Override{Exp: &Token{Token: "hello"}}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	result = trees.Fold(result)
	assert.Equal(t, "hello", asjson.AsJSON(result))
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
	m := asjson.AsJSON(result).(map[string]any)
	assert.Equal(t, "test", m["__class__"])
	assert.Equal(t, "hello", m["ast"])
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
	m := asjson.AsJSON(result).(map[string]any)
	assert.Equal(t, "start", m["__class__"])
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
	m := asjson.AsJSON(result).(map[string]any)
	assert.Equal(t, "first", m["__class__"])
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
	assert.Equal(t, "x", asjson.AsJSON(result))
}

func TestParseVoid(t *testing.T) {
	ctx := ctxFrom("hello")
	expr := &Void{}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.True(t, result == nil, "expected nil, got %T", result)
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
	assert.True(t, result == nil, "expected nil, got %T", result)
}

func TestParseConstant(t *testing.T) {
	ctx := ctxFrom("anything")
	expr := &Constant{Literal: "test"}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test", asjson.AsJSON(result))
}

func TestParseChoiceResetsCursor(t *testing.T) {
	ctx := ctxFrom("hello world")
	expr := &Choice{
		Options: []Model{
			&Option{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "wrong"}, &Token{Token: "stuff"}},
			}},
			&Option{Exp: &Sequence{
				Sequence: []Model{&Token{Token: "hello"}, &Token{Token: "world"}},
			}},
		},
	}
	result, err := expr.Parse(ctx)
	assert.NoError(t, err)
	lst := asjson.AsJSON(result).([]any)
	assert.Equal(t, 2, len(lst))
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
	lst := asjson.AsJSON(result).([]any)
	assert.Equal(t, 2, len(lst), "expected 2 items (a + closure containing b c), got %d", len(lst))
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
	result = trees.Fold(result)
	ast := asjson.AsJSON(result).(map[string]any)
	assert.Equal(t, 2, len(ast))
	assert.Equal(t, "hello", ast["first"])
	assert.Equal(t, "world", ast["second"])
}
