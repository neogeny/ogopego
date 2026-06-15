package context

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/trees"
)

func newTestBaseCtx() *CoreCtx {
	c := input.NewStrCursor("some input text")
	return NewCtx(c, nil)
}

func TestBaseCtxCallStack(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Enter("rule_a")
	ctx.Enter("rule_b")
	assert.Equal(t, 2, len(ctx.CallStack()), "expected 2 entries")
	assert.Equal(t, "rule_a", ctx.CallStack()[0])
	assert.Equal(t, "rule_b", ctx.CallStack()[1])
	ctx.Leave()
	assert.Equal(t, 1, len(ctx.CallStack()), "expected 1 entry after leave")
	assert.Equal(t, "rule_a", ctx.CallStack()[0])
	ctx.Leave()
	assert.Equal(t, 0, len(ctx.CallStack()), "expected empty callstack")
}

func TestBaseCtxCallStackLeaveEmpty(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Leave()
	assert.Equal(t, 0, len(ctx.CallStack()), "expected empty callstack after leave on empty")
}

func TestBaseCtxFailureNoRegress(t *testing.T) {
	ctx := newTestBaseCtx()
	_ = ctx.Failure(10, errors.New("first"))
	_ = ctx.Failure(3, errors.New("second"))
	furthest := ctx.FurthestFailure()
	assert.Equal(t, 10, furthest.Mark(), "expected furthest to stay 'first'")
}

func TestBaseCtxFailureProgression(t *testing.T) {
	ctx := newTestBaseCtx()
	_ = ctx.Failure(5, errors.New("first"))
	_ = ctx.Failure(15, errors.New("second"))
	furthest := ctx.FurthestFailure()
	assert.Equal(t, 15, furthest.Mark(), "expected furthest to update to 'second'")
}

func TestBaseCtxTokenMatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	assert.True(t, ctx.MatchToken("some"), `expecting: "error"`)
}

func TestBaseCtxTokenMismatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	assert.False(t, ctx.MatchToken("wrong"), "expected no match")
}

func TestBaseCtxPatternMatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	matched, err := ctx.MatchPattern(`\w+`)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, "some", matched)
}

func TestBaseCtxPatternMismatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	_, err := ctx.MatchPattern(`\d+`)
	assert.Error(t, err, "expected error")
}

func TestBaseCtxConstant(t *testing.T) {
	ctx := newTestBaseCtx()
	t1, _ := ctx.Constant("hello")
	tt, ok := t1.(*trees.Text)
	assert.True(t, ok, "expected Text, got %T", t1)
	assert.Equal(t, "hello", tt.Value)
	t2, _ := ctx.Constant(42)
	tn, ok := t2.(*trees.Number)
	assert.True(t, ok, "expected Number, got %T", t2)
	assert.Equal(t, 42, tn.Value)
	t3, _ := ctx.Constant(true)
	tb, ok := t3.(*trees.Bool)
	assert.True(t, ok, "expected Bool, got %T", t3)
	assert.Equal(t, true, tb.Value)
	t4, _ := ctx.Constant(nil)
	assert.True(t, t4 == nil, "expected nil, got %T", t4)
}

func TestBaseCtxEof(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.False(t, ctx.Eof(), "expected Eof false")
	c := input.NewStrCursor("")
	ctx2 := NewCtx(c, nil)
	assert.True(t, ctx2.Eof(), "expected Eof true on empty input")
}

func TestBaseCtxKeywords(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.False(t, ctx.IsKeyword("if"), "expected 'if' not a keyword initially")
	ctx.setKeywords([]string{"if", "else", "for"})
	assert.True(t, ctx.IsKeyword("if"), "expected 'if' to be a keyword")
	assert.True(t, ctx.IsKeyword("else"), "expected 'else' to be a keyword")
	assert.False(t, ctx.IsKeyword("while"), "expected 'while' not to be a keyword")
}

func TestBaseCtxIntern(t *testing.T) {
	ctx := newTestBaseCtx()
	s1 := ctx.Intern("hello")
	s2 := ctx.Intern("hello")
	assert.Equal(t, s2, s1, "expected interned strings to match")
}

func TestBaseCtxMarkReset(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.Equal(t, 0, ctx.Mark(), "expected Mark 0")
	ctx.Reset(5)
	assert.Equal(t, 5, ctx.Mark(), "expected Mark 5 after Reset")
}

func TestBaseCtxCursorDelegation(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.False(t, ctx.AtEnd(), "expected not at end")
	r, ok := ctx.Next()
	assert.True(t, ok, "expected ok from Next")
	assert.Equal(t, 's', r)
	r, ok = ctx.Peek()
	assert.True(t, ok, "expected ok from Peek")
	assert.Equal(t, 'o', r)
}

func TestBaseCtxDot(t *testing.T) {
	ctx := newTestBaseCtx()
	r, err := ctx.Dot()
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, 's', r)
}

func TestBaseCtxParseEOF(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.False(t, ctx.ParseEOF(), "expected ParseEOF false on non-empty input")
}

func TestBaseCtxParseEOFAtEnd(t *testing.T) {
	c := input.NewStrCursor("")
	ctx := NewCtx(c, nil)
	assert.True(t, ctx.ParseEOF(), "expected ParseEOF true on empty input")
}

func TestBaseCtxMatchToken(t *testing.T) {
	ctx := newTestBaseCtx()
	ok := ctx.MatchToken("some")
	assert.True(t, ok, "expected MatchToken 'some'")
}

func TestBaseCtxMatchPattern(t *testing.T) {
	ctx := newTestBaseCtx()
	m, err := ctx.MatchPattern(`\w+`)
	assert.NoError(t, err, "expected pattern matchi")
	assert.Equal(t, "some", m)
}

func TestBaseCtxMatchEOL(t *testing.T) {
	ctx := newTestBaseCtx()
	assert.False(t, ctx.MatchEOL(), "expected no EOL match in 'some input text'")
}

func TestBaseCtxVoid(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Void()
}

func TestBaseCtxFail(t *testing.T) {
	ctx := newTestBaseCtx()
	err := ctx.Fail()
	assert.Error(t, err, "expected error from Fail()")
}

func TestBaseCtxEofCheck(t *testing.T) {
	ctx := newTestBaseCtx()
	err := ctx.EofCheck()
	assert.Error(t, err, "expected error from EofCheck on non-empty input")
}

func TestBaseCtxEofCheckAtEnd(t *testing.T) {
	c := input.NewStrCursor("")
	ctx := NewCtx(c, nil)
	err := ctx.EofCheck()
	assert.NoError(t, err, "expected nil")
}
