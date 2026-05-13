package input

import (
	"fmt"
	"strings"
)

const (
	ansiRed   = "\033[31m"
	ansiBlue  = "\033[34m"
	ansiBold  = "\033[1m"
	ansiReset = "\033[0m"
	ansiGrey  = "\033[90m"
)

type Memento struct {
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

func NewMemento(start int, msg string, cursor Cursor, callstack []string) *Memento {
	line, col := cursor.PosAt(cursor.Mark())
	return &Memento{
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
	b.WriteString(fmt.Sprintf("    %s\n", bluePipe))

	lines := strings.Split(m.Text, "\n")
	markLineIdx := m.Line
	if markLineIdx >= len(lines) {
		markLineIdx = len(lines) - 1
	}
	startLineIdx := 0
	if markLineIdx > 4 {
		startLineIdx = markLineIdx - 4
	}

	for i := startLineIdx; i <= markLineIdx; i++ {
		b.WriteString(fmt.Sprintf("%s%2d%s %s %s\n",
			ansiBlue+ansiBold, i+1, ansiReset, bluePipe, lines[i]))

		if i == m.Line {
			pad := strings.Repeat(" ", m.Col)
			_, _ = fmt.Fprintf(&b, "    %s %s%s^%s %s%s%s\n",
				bluePipe, pad,
				ansiBold+ansiRed, ansiReset,
				ansiRed, m.Msg, ansiReset)
		}
	}

	if len(m.CallStack) > 0 {
		b.WriteString("\n")
		for _, call := range m.CallStack {
			b.WriteString(fmt.Sprintf(" %s%s→%s %s%s%s\n",
				ansiRed, ansiBold, ansiReset,
				ansiGrey, call, ansiReset))
		}
	}

	return b.String()
}

func (m *Memento) String() string {
	return m.Error()
}
