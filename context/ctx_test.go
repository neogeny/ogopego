package context

import (
	"testing"

	"github.com/neogeny/ogopego/input"
)

func newTestBaseCtx() *BaseCtx {
	c := input.NewStrCursor("some input text")
	return &BaseCtx{
		cursor:    c,
		cursorMut: c,
	}
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
	dis := ctx.Failure(5, ParseFailure{Message: "expected number"})
	if dis == nil {
		t.Fatal("expected non-nil DisasterReport")
	}
	if dis.Start != 5 {
		t.Errorf("expected Start 5, got %d", dis.Start)
	}
	if dis.Failure.Message != "expected number" {
		t.Errorf("expected 'expected number', got %q", dis.Failure.Message)
	}
	if dis.CutSeen {
		t.Error("expected CutSeen false")
	}
	if ctx.FurthestFailure() != dis {
		t.Error("expected furthest failure to match")
	}
}

func TestBaseCtxFailureNoRegress(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Failure(10, ParseFailure{Message: "first"})
	ctx.Failure(3, ParseFailure{Message: "second"})
	if ctx.FurthestFailure().Failure.Message != "first" {
		t.Errorf("expected furthest to stay 'first', got %q", ctx.FurthestFailure().Failure.Message)
	}
}

func TestBaseCtxFailureProgression(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Failure(5, ParseFailure{Message: "first"})
	ctx.Failure(15, ParseFailure{Message: "second"})
	if ctx.FurthestFailure().Failure.Message != "second" {
		t.Errorf("expected furthest 'second', got %q", ctx.FurthestFailure().Failure.Message)
	}
}

func TestBaseCtxCut(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.CutSeen() {
		t.Error("expected CutSeen false initially")
	}
	ctx.Cut()
	if !ctx.CutSeen() {
		t.Error("expected CutSeen true after Cut()")
	}
	ctx.ClearCut()
	if ctx.CutSeen() {
		t.Error("expected CutSeen false after ClearCut()")
	}
}

func TestBaseCtxCutInFailure(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.Cut()
	dis := ctx.Failure(5, ParseFailure{Message: "cut failure"})
	if !dis.CutSeen {
		t.Error("expected CutSeen true in failure report")
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

func TestBaseCtxMatchEOL(t *testing.T) {
	ctx := newTestBaseCtx()
	if ctx.MatchEOL() {
		t.Error("expected no EOL match in 'some input text'")
	}
}

func TestBaseCtxMatchToken(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	if !ctx.MatchToken("some") {
		t.Error("expected MatchToken 'some'")
	}
}

func TestBaseCtxMatchPattern(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.NextToken()
	m, ok := ctx.MatchPattern(`\w+`)
	if !ok {
		t.Fatal("expected pattern match")
	}
	if m != "some" {
		t.Errorf("expected 'some', got %q", m)
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
	ctx := &BaseCtx{
		cursor:    c,
		cursorMut: c,
	}
	if !ctx.ParseEOF() {
		t.Error("expected ParseEOF true on empty input")
	}
}

func TestBaseCtxPushMerge(t *testing.T) {
	ctx := newTestBaseCtx()
	if p := ctx.Push(); p != nil {
		t.Error("expected nil from Push")
	}
	if m := ctx.Merge(nil); m != nil {
		t.Error("expected nil from Merge")
	}
}

func TestBaseCtxMemo(t *testing.T) {
	ctx := newTestBaseCtx()
	key := MemoKey{Mark: 0, Name: "test", CanMemo: true}
	if m := ctx.Memo(key); m != nil {
		t.Error("expected nil memo")
	}
	ctx.Memoize(key, "some tree", 10)
}

func TestBaseCtxClearErrorMemos(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.ClearErrorMemos()
}

func TestBaseCtxPruneCache(t *testing.T) {
	ctx := newTestBaseCtx()
	ctx.PruneCache(5)
}

func TestBaseCtxTrackUntrack(t *testing.T) {
	ctx := newTestBaseCtx()
	key := MemoKey{Mark: 0, Name: "test", CanMemo: true}
	if n := ctx.Track(key); n != 0 {
		t.Errorf("expected 0 from Track, got %d", n)
	}
	if n := ctx.Untrack(key); n != 0 {
		t.Errorf("expected 0 from Untrack, got %d", n)
	}
}
