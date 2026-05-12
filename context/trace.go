package context

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Event int

const (
	EventEntry    Event = iota
	EventSuccess
	EventFailure
	EventRecursion
	EventCut
	EventMatch
	EventNoMatch
)

type Tracer interface {
	Trace(ctx Ctx, msg string)
	TraceEvent(ctx Ctx, event Event, msg string)
	TraceEntry(ctx Ctx)
	TraceSuccess(ctx Ctx)
	TraceFailure(ctx Ctx, err string)
	TraceRecursion(ctx Ctx)
	TraceCut(ctx Ctx)
	TraceMatch(ctx Ctx, token, name string) bool
	TraceNoMatch(ctx Ctx, token, name string) bool
}

type NullTracer struct{}

func (NullTracer) Trace(_ Ctx, _ string)                           {}
func (NullTracer) TraceEvent(_ Ctx, _ Event, _ string)              {}
func (NullTracer) TraceEntry(_ Ctx)                                 {}
func (NullTracer) TraceSuccess(_ Ctx)                               {}
func (NullTracer) TraceFailure(_ Ctx, _ string)                     {}
func (NullTracer) TraceRecursion(_ Ctx)                             {}
func (NullTracer) TraceCut(_ Ctx)                                   {}
func (NullTracer) TraceMatch(_ Ctx, _, _ string) bool               { return true }
func (NullTracer) TraceNoMatch(_ Ctx, _, _ string) bool             { return false }

type ConsoleTracer struct{}

func eventSymbol(event Event) string {
	switch event {
	case EventEntry:
		return color.YellowString("↙")
	case EventSuccess:
		return color.GreenString("≡")
	case EventFailure:
		return color.RedString("≢")
	case EventRecursion:
		return color.BlueString("⟲")
	case EventCut:
		return color.YellowString("⚔")
	case EventMatch:
		return color.GreenString("≡")
	case EventNoMatch:
		return color.RedString("≢")
	default:
		return "?"
	}
}

func stackSymbol(event Event) string {
	switch event {
	case EventSuccess:
		return color.GreenString("→")
	case EventFailure:
		return color.RedString("→")
	case EventNoMatch:
		return color.RedString("←")
	case EventMatch:
		return color.GreenString("←")
	default:
		return color.YellowString("←")
	}
}

func (ConsoleTracer) Trace(_ Ctx, msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func (t ConsoleTracer) TraceEvent(ctx Ctx, event Event, msg string) {
	esym := eventSymbol(event)
	ssym := stackSymbol(event)

	lookahead := color.New(color.FgBlack, color.Faint).Sprintf(
		"%s",
		ctx.Cursor().Lookahead(ctx.Mark()),
	)

	var cs string
	for _, call := range ctx.CallStack() {
		cs += color.WhiteString(call) + ssym
	}

	line, col := ctx.Cursor().Pos()
	pos := color.New(color.FgBlack, color.Faint).Sprintf("[%d:%d]→", line, col)

	lineMsg := fmt.Sprintf("%s%s %s •\n%s%s",
		esym, msg, cs, pos, lookahead)

	t.Trace(ctx, lineMsg)
}

func (t ConsoleTracer) TraceEntry(ctx Ctx) {
	t.TraceEvent(ctx, EventEntry, "")
}

func (t ConsoleTracer) TraceSuccess(ctx Ctx) {
	t.TraceEvent(ctx, EventSuccess, "")
}

func (t ConsoleTracer) TraceFailure(ctx Ctx, err string) {
	errStr := fmt.Sprintf(" %s", color.RedString(err))
	t.TraceEvent(ctx, EventFailure, errStr)
}

func (t ConsoleTracer) TraceRecursion(ctx Ctx) {
	t.TraceEvent(ctx, EventRecursion, "")
}

func (t ConsoleTracer) TraceCut(ctx Ctx) {
	t.TraceEvent(ctx, EventCut, "")
}

func (t ConsoleTracer) TraceMatch(ctx Ctx, token, name string) bool {
	tag := ""
	if name != "" {
		tag = fmt.Sprintf("/%s/", name)
	}
	msg := color.GreenString("'%s'%s", token, tag)
	t.TraceEvent(ctx, EventMatch, msg)
	return true
}

func (t ConsoleTracer) TraceNoMatch(ctx Ctx, token, name string) bool {
	var msg string
	if token != "" {
		msg = color.RedString(" '%s'", token)
	} else {
		msg = color.RedString(" /%s/", name)
	}
	t.TraceEvent(ctx, EventNoMatch, msg)
	return false
}
