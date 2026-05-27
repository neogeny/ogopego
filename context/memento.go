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
			"  %s %s @ [%d:%d]\n",
			arrow,
			m.Cursor.InputSource(),
			line,
			col,
		),
	)

	b.WriteString(fmt.Sprintf(" %5s%s\n", "", bluePipe))
	start := line - 4
	if start < 0 {
		start = 0
	}
	for i, linestr := range m.Cursor.LinesAt(start, line+1) {
		linestr = util.StripRight(linestr)
		disp := util.ExpandTabs(linestr)
		b.WriteString(
			fmt.Sprintf(
				"%s%5d %s %s\n",
				ansiBlue+ansiBold,
				start+i,
				ansiReset+bluePipe,
				disp,
			),
		)
	}
	pad := strings.Repeat(" ", col)
	_, _ = fmt.Fprintf(&b, " %5s%s %s^ %s\n",
		"",
		bluePipe,
		ansiReset+pad+ansiRed,
		m.Msg+ansiReset,
	)

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
