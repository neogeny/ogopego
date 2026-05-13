package context

import (
	"errors"
	"testing"

	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/trees"
)

func newTestBaseCtx() *BaseCtx {
	c := input.NewStrCursor("some input text")
	return NewBaseCtx(c)
}

func TestBaseCtxCallStack(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Enter("rule_a")
	ctx.Enter("rule_b")
	if len(ctx.CallStack()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ctx.CallStack()))
	}
	if ctx.CallStack()[0] != "rule_a" {
		t.Errorf("expected 'rule_a', got %q", ctx.CallStack()[0])
	}
	if ctx.CallStack()[1] != "rule_b" {
		t.Errorf("expected 'rule_b', got %q", ctx.CallStack()[1])
	}
	ctx.Leave()
	if len(ctx.CallStack()) != 1 {
		t.Fatalf("expected 1 entry after leave, got %d", len(ctx.CallStack()))
	}
	if ctx.CallStack()[0] != "rule_a" {
		t.Errorf("expected 'rule_a', got %q", ctx.CallStack()[0])
	}
	ctx.Leave()
	if len(ctx.CallStack()) != 0 {
		t.Errorf("expected empty callstack, got %d", len(ctx.CallStack()))
	}
}

func TestBaseCtxCallStackLeaveEmpty(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Leave()
	if len(ctx.CallStack()) != 0 {
		t.Errorf("expected empty callstack after leave on empty, got %d", len(ctx.CallStack()))
	}
}

func TestBaseCtxFailure(t *testing.T) {
	ctx := newTestBaseCtx()
	fail := ctx.Failure(5, errors.New("expected number"))
	if fail == nil {
		t.Fatal("expected non-nil DisasterReport")
	}
	if fail.Start != 5 {
		t.Errorf("expected Start 5, got %d", fail.Start)
	}
	if fail.Inner.Error() != "expected number" {
		t.Errorf("expected 'expected number', got %q", fail.Inner.Error())
	}
	if ctx.FurthestFailure().Failure.Inner != fail.Inner {
		t.Error("expected furthest failure to match")
	}
}

func TestBaseCtxFailureNoRegress(t *testing.T) {
	ctx := newTestBaseCtx()
	_ = ctx.Failure(10, errors.New("first"))
	_ = ctx.Failure(3, errors.New("second"))
	if ctx.FurthestFailure().Failure.Error() != "at 10: first" {
		t.Errorf("expected furthest to stay 'first', got %q", ctx.FurthestFailure().Failure.Error())
	}
}

func TestBaseCtxFailureProgression(t *testing.T) {
	ctx := newTestBaseCtx()
	_ = ctx.Failure(5, errors.New("first"))
	_ = ctx.Failure(15, errors.New("second"))
	if ctx.FurthestFailure().Failure.Error() != "at 15: second" {
		t.Errorf("expected furthest 'second', got %q", ctx.FurthestFailure().Failure.Error())
	}
}

func TestBaseCtxTokenMatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	matched, err := ctx.Token("some")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched != "some" {
		t.Errorf("expected 'some', got %q", matched)
	}
}

func TestBaseCtxTokenMismatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	_, err := ctx.Token("wrong")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBaseCtxPatternMatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	matched, err := ctx.Pattern(`\w+`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched != "some" {
		t.Errorf("expected 'some', got %q", matched)
	}
}

func TestBaseCtxPatternMismatch(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	_, err := ctx.Pattern(`\d+`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBaseCtxConstant(t *testing.T) {
	ctx := newTestBaseCtx()
	t1, _ := ctx.Constant("hello")
	if tt, ok := t1.(*trees.Text); !ok || tt.Value != "hello" {
		t.Errorf("expected Text{hello}, got %T %+v", t1, t1)
	}
	t2, _ := ctx.Constant(42)
	if tn, ok := t2.(*trees.Number); !ok || tn.Value != 42 {
		t.Errorf("expected Number{42}, got %T %+v", t2, t2)
	}
	t3, _ := ctx.Constant(true)
	if tb, ok := t3.(*trees.Bool); !ok || tb.Value != true {
		t.Errorf("expected Bool{true}, got %T %+v", t3, t3)
	}
	t4, _ := ctx.Constant(nil)
	if _, ok := t4.(*trees.Nil); !ok {
		t.Errorf("expected Nil, got %T", t4)
	}
}

func TestBaseCtxEof(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.Eof() {
		t.Error("expected Eof false")
	}
	c := input.NewStrCursor("")
	ctx2 := NewBaseCtx(c)
	if !ctx2.Eof() {
		t.Error("expected Eof true on empty input")
	}
}

func TestBaseCtxKeywords(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.IsKeyword("if") {
		t.Error("expected 'if' not a keyword initially")
	}
	ctx.SetKeywords([]string{"if", "else", "for"})
	if !ctx.IsKeyword("if") {
		t.Error("expected 'if' to be a keyword")
	}
	if !ctx.IsKeyword("else") {
		t.Error("expected 'else' to be a keyword")
	}
	if ctx.IsKeyword("while") {
		t.Error("expected 'while' not to be a keyword")
	}
}

func TestBaseCtxIntern(t *testing.T) {
	ctx := newTestBaseCtx()
	s1 := ctx.Intern("hello")
	s2 := ctx.Intern("hello")
	if s1 != s2 {
		t.Error("expected interned strings to match")
	}
}

func TestBaseCtxGetPattern(t *testing.T) {
	ctx := newTestBaseCtx()
	p1 := ctx.GetPattern(`\d+`)
	if p1 == nil {
		t.Fatal("expected non-nil pattern")
	}
	p2 := ctx.GetPattern(`\d+`)
	if p2 != p1 {
		t.Error("expected cached pattern to be same instance")
	}
}

func TestBaseCtxGetPatternInvalid(t *testing.T) {
	ctx := newTestBaseCtx()
	p := ctx.GetPattern(`[invalid`)
	if p != nil {
		t.Error("expected nil for invalid pattern")
	}
}

func TestBaseCtxMarkReset(t *testing.T) {
	ctx := newTestBaseCtx()
	if m := ctx.Mark(); m != 0 {
		t.Errorf("expected Mark 0, got %d", m)
	}
	ctx.Reset(5)
	if m := ctx.Mark(); m != 5 {
		t.Errorf("expected Mark 5 after Reset, got %d", m)
	}
}

func TestBaseCtxCursorDelegation(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.AtEnd() {
		t.Error("expected not at end")
	}
	r, ok := ctx.Next()
	if !ok || r != 's' {
		t.Errorf("expected 's', got %c", r)
	}
	r, ok = ctx.Peek()
	if !ok || r != 'o' {
		t.Errorf("expected 'o', got %c", r)
	}
}

func TestBaseCtxDot(t *testing.T) {
	ctx := newTestBaseCtx()
	r, err := ctx.Dot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != 's' {
		t.Errorf("expected 's', got %c", r)
	}
}

func TestBaseCtxParseEOF(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.ParseEOF() {
		t.Error("expected ParseEOF false on non-empty input")
	}
}

func TestBaseCtxParseEOFAtEnd(t *testing.T) {
	c := input.NewStrCursor("")
	ctx := NewBaseCtx(c)
	if !ctx.ParseEOF() {
		t.Error("expected ParseEOF true on empty input")
	}
}

func TestBaseCtxMatchToken(t *testing.T) {
	ctx := newTestBaseCtx()
	ok := ctx.MatchToken("some")
	if !ok {
		t.Error("expected MatchToken 'some'")
	}
}

func TestBaseCtxMatchPattern(t *testing.T) {
	ctx := newTestBaseCtx()
	m, ok := ctx.MatchPattern(`\w+`)
	if !ok {
		t.Fatal("expected pattern match")
	}
	if m != "some" {
		t.Errorf("expected 'some', got %q", m)
	}
}

func TestBaseCtxMatchEOL(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.MatchEOL() {
		t.Error("expected no EOL match in 'some input text'")
	}
}

func TestBaseCtxVoid(t *testing.T) {
	ctx := newTestBaseCtx()
	err := ctx.Void()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestBaseCtxFail(t *testing.T) {
	ctx := newTestBaseCtx()
	err := ctx.Fail()
	if err == nil {
		t.Error("expected error from Fail()")
	}
}

func TestBaseCtxEofCheck(t *testing.T) {
	ctx := newTestBaseCtx()
	err := ctx.EofCheck()
	if err == nil {
		t.Error("expected error from EofCheck on non-empty input")
	}
}

func TestBaseCtxEofCheckAtEnd(t *testing.T) {
	c := input.NewStrCursor("")
	ctx := NewBaseCtx(c)
	err := ctx.EofCheck()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
