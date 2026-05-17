package context

import (
	"fmt"
	"strings"

	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/util"
)

const (
	ansiRed   = "\033[31m"
	ansiBlue  = "\033[34m"
	ansiBold  = "\033[1m"
	ansiReset = "\033[0m"
	ansiGrey  = "\033[90m"
)

// Memento captures the state of the parser at a specific point for error reporting.
type Memento struct {
	Cursor    input.Cursor
	Msg       string
	Start     int
	Mark      int
	CallStack []string
}

// NewMemento constructs a Memento capturing cursor state and a message for
// later diagnostic reporting.
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

// InputSource returns the name of the input source.
func (m *Memento) InputSource() string {
	return m.Cursor.InputSource()
}

// Text returns the full text of the input.
func (m *Memento) Text() string {
	return m.Cursor.AsStr()
}

// Error returns a formatted string representation of the Memento, suitable for error messages.
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
			"  %s %s@[%d:%d]\n",
			arrow,
			m.Cursor.InputSource(),
			line,
			col,
		),
	)

	b.WriteString(fmt.Sprintf("   %s\n", bluePipe))
	i := 0
	for linestr := range strings.Lines(m.Cursor.AsStr()) {
		if i > line {
			break
		}
		linestr = strings.TrimRight(linestr, "\n\r\t\f")
		if i >= line-4 {
			disp := util.ExpandTabs(linestr)
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
		i += 1
	}
	pad := strings.Repeat(" ", col)
	_, _ = fmt.Fprintf(&b, "   %s %s%s^%s %s%s%s\n",
		bluePipe, pad,
		ansiBold+ansiRed, ansiReset,
		ansiRed, m.Msg, ansiReset)

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

// String returns a formatted string representation of the Memento.
func (m *Memento) String() string {
	return m.Error()
}
