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
	Cursor    input.Cursor
	Msg       string
	Start     int
	Mark      int
	CallStack []string
}

func NewMemento(start int, msg string, cursor input.Cursor, callstack []string) Memento {
	cs := make([]string, len(callstack))
	copy(cs, callstack)
	return Memento{
		Cursor:    cursor,
		Msg:       msg,
		Start:     start,
		Mark:      cursor.Mark(),
		CallStack: cs,
	}
}

func (m *Memento) InputSource() string {
	return m.Cursor.InputSource()
}

func (m *Memento) Text() string {
	return m.Cursor.AsStr()
}

func (m *Memento) Error() string {
	line, col := m.Cursor.PosAt(m.Mark)
	var b strings.Builder

	errLabel := ansiBold + ansiRed + "error" + ansiReset
	bluePipe := ansiBlue + ansiBold + "|" + ansiReset
	arrow := ansiBlue + ansiBold + "-->" + ansiReset

	b.WriteString("\n")
	b.WriteString(
		fmt.Sprintf(
			"%s: %s%s%s\n",
			errLabel,
			ansiBold,
			m.Msg,
			ansiReset,
		),
	)
	b.WriteString(
		fmt.Sprintf(
			"  %s %s:%d:%d\n",
			arrow,
			m.Cursor.InputSource(),
			line,
			col,
		),
	)

	b.WriteString(fmt.Sprintf("   %s\n", bluePipe))
	i := 1
	for linestr := range strings.Lines(m.Cursor.AsStr()) {
		if i > line {
			break
		}
		linestr = strings.TrimRight(linestr, "\n\r\t\f")
		if i >= line-4 {
			disp := expandTabs(linestr)
			b.WriteString(
				fmt.Sprintf(
					"%s%2d%s %s %s\n",
					ansiBlue+ansiBold,
					i,
					ansiReset,
					bluePipe,
					disp,
				),
			)
		}
		if i == line {
			pad := strings.Repeat(" ", col)
			_, _ = fmt.Fprintf(&b, "   %s %s%s^%s %s%s%s\n",
				bluePipe, pad,
				ansiBold+ansiRed, ansiReset,
				ansiRed, m.Msg, ansiReset)
		}
		i += 1
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
