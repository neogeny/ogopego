package context

import (
	"strings"
	"testing"

	"github.com/neogeny/ogopego/input"
)

func newTestCtx() *BaseCtx {
	c := input.NewStrCursor("some input text")
	return NewBaseCtx(c)
}

func newTestCtxWithTracer(tracer Tracer) *BaseCtx {
	c := input.NewStrCursor("some input text")
	return NewBaseCtxWithTracer(c, tracer)
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
	if !result {
		t.Error("TraceMatch should return true")
	}
}

func TestConsoleTracerTraceNoMatch(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceNoMatch(ctx, "", "pattern")
	if result {
		t.Error("TraceNoMatch should return false")
	}
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
	if EventEntry != 0 {
		t.Errorf("EventEntry should be 0, got %d", EventEntry)
	}
	if EventNoMatch != 6 {
		t.Errorf("EventNoMatch should be 6, got %d", EventNoMatch)
	}
}

func TestBaseCtxTracerDefault(t *testing.T) {
	ctx := newTestCtx()
	tr := ctx.Tracer()
	if tr == nil {
		t.Error("expected non-nil default tracer")
	}
}

func TestBaseCtxWithTracer(t *testing.T) {
	ct := ConsoleTracer{}
	ctx := newTestCtxWithTracer(ct)
	if ctx.Tracer() != ct {
		t.Error("expected tracer to match")
	}
}

func TestTraceWithCallstack(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	ctx.Enter("grammar")
	ctx.Enter("expr")
	ctx.Enter("term")
	et := ctx.Tracer()
	et.TraceEntry(ctx)
	if len(ctx.CallStack()) != 3 {
		t.Errorf("expected 3 callstack entries, got %d", len(ctx.CallStack()))
	}
}

func TestBaseCtxNewWithTracer(t *testing.T) {
	c := input.NewStrCursor("test")
	ctx := NewBaseCtxWithTracer(c, ConsoleTracer{})
	_, ok := ctx.Tracer().(ConsoleTracer)
	if !ok {
		t.Error("expected ConsoleTracer")
	}
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
	if !result {
		t.Error("TraceMatch should return true")
	}
}

func TestConsoleTracerTraceNoMatchPattern(t *testing.T) {
	ctx := newTestCtxWithTracer(ConsoleTracer{})
	et := ctx.Tracer()
	result := et.TraceNoMatch(ctx, "", "\\d+")
	if result {
		t.Error("TraceNoMatch should return false")
	}
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
		if !strings.Contains(got, tt.want) {
			t.Errorf("eventSymbol(%d) should contain %q, got %q", tt.event, tt.want, got)
		}
	}
}
