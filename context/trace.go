package context

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/util"
)

// Event represents a tracing event kind emitted during parsing.
type Event int

// Event enumerates tracing event kinds emitted during parsing.
//
// Values include entry, success, failure, recursion, cut, match and no-match.
// These are supplied to TraceEvent by Tracer implementations.
const (
	EventEntry Event = iota
	EventSuccess
	EventFailure
	EventRecursion
	EventCut
	EventMatch
	EventNoMatch
)

// Tracer defines the interface for tracing parsing events.
type Tracer interface {
	// Trace logs a general message.
	Trace(ctx Ctx, msg string)
	// TraceEvent logs a specific parsing event with a message.
	TraceEvent(ctx Ctx, event Event, msg string)
	// TraceEntry logs the entry into a parsing rule.
	TraceEntry(ctx Ctx)
	// TraceSuccess logs a successful parsing rule.
	TraceSuccess(ctx Ctx)
	// TraceFailure logs a failed parsing rule with an error message.
	TraceFailure(ctx Ctx, err string)
	// TraceRecursion logs a recursive call to a parsing rule.
	TraceRecursion(ctx Ctx)
	// TraceCut logs a cut operation.
	TraceCut(ctx Ctx)
	// TraceMatch logs a successful token match.
	TraceMatch(ctx Ctx, token, name string) bool
	// TraceNoMatch logs a failed token match.
	TraceNoMatch(ctx Ctx, token, name string) bool
}

// NullTracer is a no-op tracer used when tracing is disabled.
type NullTracer struct{}

// Trace does nothing for NullTracer.
func (NullTracer) Trace(_ Ctx, _ string) {}

// TraceEvent does nothing for NullTracer.
func (NullTracer) TraceEvent(_ Ctx, _ Event, _ string) {}

// TraceEntry does nothing for NullTracer.
func (NullTracer) TraceEntry(_ Ctx) {}

// TraceSuccess does nothing for NullTracer.
func (NullTracer) TraceSuccess(_ Ctx) {}

// TraceFailure does nothing for NullTracer.
func (NullTracer) TraceFailure(_ Ctx, _ string) {}

// TraceRecursion does nothing for NullTracer.
func (NullTracer) TraceRecursion(_ Ctx) {}

// TraceCut does nothing for NullTracer.
func (NullTracer) TraceCut(_ Ctx) {}

// TraceMatch returns true for NullTracer.
func (NullTracer) TraceMatch(_ Ctx, _, _ string) bool { return true }

// TraceNoMatch returns false for NullTracer.
func (NullTracer) TraceNoMatch(_ Ctx, _, _ string) bool { return false }

// ConsoleTracer implements Tracer to output tracing information to the console.
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

// Trace prints a message to stderr.
func (ConsoleTracer) Trace(_ Ctx, msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// TraceEvent formats and prints a tracing event to stderr.
func (t ConsoleTracer) TraceEvent(ctx Ctx, event Event, msg string) {
	esym := eventSymbol(event)
	ssym := stackSymbol(event)

	lookahead := color.New(color.FgBlack, color.Bold).Sprintf(
		"%s",
		strings.ReplaceAll(ctx.Cursor().Lookahead(ctx.Mark()), " ", "·"),
	)

	cols, _ := util.TermSize()
	var cs string
	bold := color.New(color.FgWhite, color.Bold)
	callStack := ctx.CallStack()
	for i := len(callStack) - 1; i >= 0; i-- {
		cs += bold.Sprint(callStack[i]) + ssym
		if len(cs) >= cols-8 {
			cs += "••"
			break
		}
	}

	line, col := ctx.Cursor().Pos()
	pos := color.New(color.FgBlack, color.Bold).Sprintf("[%d:%d]→", line, col)

	lineMsg := fmt.Sprintf("%s%s %s•\n%s%s",
		esym, msg, cs, pos, lookahead)

	t.Trace(ctx, lineMsg)
}

// TraceEntry logs the entry into a parsing rule.
func (t ConsoleTracer) TraceEntry(ctx Ctx) {
	t.TraceEvent(ctx, EventEntry, "")
}

// TraceSuccess logs a successful parsing rule.
func (t ConsoleTracer) TraceSuccess(ctx Ctx) {
	t.TraceEvent(ctx, EventSuccess, "")
}

// TraceFailure logs a failed parsing rule with an error message.
func (t ConsoleTracer) TraceFailure(ctx Ctx, err string) {
	errStr := fmt.Sprintf(" %s", color.RedString(err))
	t.TraceEvent(ctx, EventFailure, errStr)
}

// TraceRecursion logs a recursive call to a parsing rule.
func (t ConsoleTracer) TraceRecursion(ctx Ctx) {
	t.TraceEvent(ctx, EventRecursion, "")
}

// TraceCut logs a cut operation.
func (t ConsoleTracer) TraceCut(ctx Ctx) {
	t.TraceEvent(ctx, EventCut, "")
}

// TraceMatch logs a successful token match.
func (t ConsoleTracer) TraceMatch(ctx Ctx, token, name string) bool {
	tag := ""
	if name != "" {
		tag = fmt.Sprintf("/%s/", name)
	}
	msg := color.GreenString("'%s'%s", token, tag)
	t.TraceEvent(ctx, EventMatch, msg)
	return true
}

// TraceNoMatch logs a failed token match.
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
