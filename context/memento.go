package context

import (
	"fmt"
	"strings"

	"github.com/neogeny/ogopego/input"
)

const (
	ansiRed   = "\033[31m"
	ansiBlue  = "\033[34m"
	ansiBold  = "\033[1m"
	ansiReset = "\033[0m"
	ansiGrey  = "\033[90m"
)

func expandTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

type Memento struct {
	error
	InputSource string
	Msg         string
	Text        string
	Start       int
	Mark        int
	CallStack   []string
	Line        int
	Col         int
	LA          string
}

func NewMemento(start int, msg string, cursor input.Cursor, callstack []string) Memento {
	line, col := cursor.PosAt(cursor.Mark())
	return Memento{
		InputSource: cursor.InputSource(),
		Msg:         msg,
		Text:        cursor.AsStr(),
		Start:       start,
		Mark:        cursor.Mark(),
		CallStack:   callstack,
		Line:        line,
		Col:         col,
		LA:          cursor.Lookahead(start),
	}
}

func (m *Memento) Error() string {
	var b strings.Builder

	errLabel := ansiBold + ansiRed + "error" + ansiReset
	bluePipe := ansiBlue + ansiBold + "|" + ansiReset
	arrow := ansiBlue + ansiBold + "-->" + ansiReset

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%s: %s%s%s\n", errLabel, ansiBold, m.Msg, ansiReset))
	b.WriteString(fmt.Sprintf("  %s %s:%d:%d\n", arrow, m.InputSource, m.Line, m.Col))

	lines := strings.Split(m.Text, "\n")
	markLineIdx := m.Line
	if markLineIdx >= len(lines) {
		markLineIdx = len(lines) - 1
	}
	startLineIdx := 0
	if markLineIdx > 4 {
		startLineIdx = markLineIdx - 4
	}

	b.WriteString(fmt.Sprintf("   %s\n", bluePipe))
	for i := startLineIdx; i <= markLineIdx; i++ {
		disp := expandTabs(lines[i])
		b.WriteString(fmt.Sprintf("%s%2d%s %s %s\n",
			ansiBlue+ansiBold, i+1, ansiReset, bluePipe, disp))

		if i == m.Line {
			pad := strings.Repeat(" ", m.Col)
			_, _ = fmt.Fprintf(&b, "   %s %s%s^%s %s%s%s\n",
				bluePipe, pad,
				ansiBold+ansiRed, ansiReset,
				ansiRed, m.Msg, ansiReset)
		}
	}

	if len(m.CallStack) > 0 {
		b.WriteString("\n")
		for i := len(m.CallStack) - 1; i >= 0; i-- {
			b.WriteString(fmt.Sprintf(" %s%s→%s %s%s%s\n",
				ansiRed, ansiBold, ansiReset,
				ansiGrey, m.CallStack[i], ansiReset))
		}
	}

	return b.String()
}

func (m *Memento) String() string {
	return m.Error()
}
