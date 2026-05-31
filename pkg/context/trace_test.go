package context

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/input"
)

func newTestCtx() *CoreCtx {
	c := input.NewStrCursor("some input text")
	return NewCtx(c, nil)
}

func newTestCtxWithTracer(tracer Tracer) *CoreCtx {
	c := input.NewStrCursor("some input text")
	ctx := NewCtx(c, nil)
	ctx.SetTracer(tracer)
	return ctx
}

func TestNullTracerImplementsTracer(t *testing.T) {
	var nt NullTracer
	var iface Tracer = nt
	_ = iface
}

func TestConsoleTracerImplementsTracer(t *testing.T) {
	var ct ConsoleTracer
	var iface Tracer = ct
	_ = iface
}

func TestNullTracerMethods(t *testing.T) {
	ctx := newTestCtx()
	nt := NullTracer{}
	nt.Trace(ctx, "msg")
	nt.TraceEvent(ctx, EventEntry, "msg")
	nt.TraceEntry(ctx)
	nt.TraceSuccess(ctx)
	nt.TraceFailure(ctx, "err")
	nt.TraceRecursion(ctx)
	nt.TraceCut(ctx)
	nt.TraceMatch(ctx, "token", "rule")
	nt.TraceNoMatch(ctx, "", "pattern")
}

func TestConsoleTracerTraceEntry(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEntry(ctx)
}

func TestConsoleTracerTraceSuccess(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceSuccess(ctx)
}

func TestConsoleTracerTraceFailure(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceFailure(ctx, "expected token")
}

func TestConsoleTracerTraceMatch(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceMatch(ctx, "hello", "greeting")
	assert.True(t, result, "TraceMatch should return true")
}

func TestConsoleTracerTraceNoMatch(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceNoMatch(ctx, "", "pattern")
	assert.False(t, result, "TraceNoMatch should return false")
}

func TestConsoleTracerTraceRecursion(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	ctx.Enter("expr")
	ctx.Enter("expr")
	et := ctx.Tracer()
	et.TraceRecursion(ctx)
}

func TestConsoleTracerTraceCut(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceCut(ctx)
}

func TestNullTracerMethodsNoPanic(t *testing.T) {
	ctx := newTestCtx()
	nt := NullTracer{}
	// Verify no panics
	nt.TraceEntry(ctx)
	nt.TraceSuccess(ctx)
	nt.TraceFailure(ctx, "err")
	nt.TraceRecursion(ctx)
	nt.TraceCut(ctx)
	nt.TraceMatch(ctx, "tok", "rule")
	nt.TraceNoMatch(ctx, "", "pattern")
}

func TestEventConstants(t *testing.T) {
	assert.Equal(t, 0, EventEntry, "EventEntry should be 0")
	assert.Equal(t, 6, EventNoMatch, "EventNoMatch should be 6")
}

func TestBaseCtxTracerDefault(t *testing.T) {
	ctx := newTestCtx()
	tr := ctx.Tracer()
	assert.NotZero(t, tr, "expected non-nil default tracer")
}

func TestBaseCtxWithTracer(t *testing.T) {
	ct := ConsoleTracer{}
	ctx := newTestCtxWithTracer(ct)
	assert.True(t, ctx.Tracer() == ct, "expected tracer to match")
}

func TestTraceWithCallstack(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	ctx.Enter("grammar")
	ctx.Enter("expr")
	ctx.Enter("term")
	et := ctx.Tracer()
	et.TraceEntry(ctx)
	assert.Equal(t, 3, len(ctx.CallStack()), "expected 3 callstack entries")
}

func TestBaseCtxNewWithTracer(t *testing.T) {
	c := input.NewStrCursor("test")
	ctx := NewCtx(c, nil)
	ctx.SetTracer(ConsoleTracer{})
	_, ok := ctx.Tracer().(ConsoleTracer)
	assert.True(t, ok, "expected ConsoleTracer")
}

func TestConsoleTracerTraceEventMatch(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEvent(ctx, EventMatch, "hello")
}

func TestConsoleTracerTraceEventNoMatch(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEvent(ctx, EventNoMatch, "pattern")
}

func TestConsoleTracerTraceEventFailure(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEvent(ctx, EventFailure, "syntax error")
}

func TestConsoleTracerTraceEventRecursion(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEvent(ctx, EventRecursion, "")
}

func TestConsoleTracerTraceEventCut(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	et.TraceEvent(ctx, EventCut, "")
}

func TestConsoleTracerTraceMatchWithName(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceMatch(ctx, "if", "keyword")
	assert.True(t, result, "TraceMatch should return true")
}

func TestConsoleTracerTraceNoMatchPattern(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceNoMatch(ctx, "", "\\d+")
	assert.False(t, result, "TraceNoMatch should return false")
}

func TestEventSymbols(t *testing.T) {
	tests := []struct {
		event Event
		want  string
	}{
		{EventEntry, "↙"},
		{EventSuccess, "≡"},
		{EventFailure, "≢"},
		{EventRecursion, "⟲"},
		{EventCut, "⚔"},
		{EventMatch, "≡"},
		{EventNoMatch, "≢"},
	}
	for _, tt := range tests {
		got := eventSymbol(tt.event)
		assert.True(t, strings.Contains(got, tt.want), "eventSymbol(%d) should contain %q, got %q", tt.event, tt.want, got)
	}
}
