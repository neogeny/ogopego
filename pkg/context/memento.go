package context

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/util"
)

var (
	errStyle    = color.New(color.FgRed, color.Bold)
	gutStyle    = color.New(color.FgBlue, color.Bold)
	bold        = color.New(color.Bold)
	grey        = color.New(color.FgHiBlack)
	dimWhite    = color.New(color.FgWhite, color.Faint)
	dimRedStyle = color.New(color.FgRed, color.Faint)
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

	errLabel := errStyle.Sprint("error")
	gut := gutStyle.Sprint("│")
	arrow := gutStyle.Sprint("─→")

	b.WriteString(fmt.Sprintf("\n%s: %s\n", errLabel, bold.Sprint(m.Msg)))
	b.WriteString(fmt.Sprintf("  %s %s[%d:%d]\n", arrow, m.Cursor.InputSource(), line, col))

	b.WriteString(fmt.Sprintf(" %5s%s\n", "", gut))
	start := line - 4
	if start < 0 {
		start = 0
	}
	for i, linestr := range m.Cursor.LinesAt(start, line+1) {
		linestr = util.StripRight(linestr)
		disp := util.ExpandTabs(linestr)
		lineno := dimWhite.Sprintf("%5d", start+i)
		b.WriteString(fmt.Sprintf("%s %s %s\n", lineno, gut, disp))
	}
	pad := strings.Repeat(" ", col-1)
	_, _ = fmt.Fprintf(&b, " %5s%s %s\n", "", gut, errStyle.Sprintf("%s⌃ %s", pad, m.Msg))

	if len(m.CallStack) > 0 {
		b.WriteString("\n")
		for i := len(m.CallStack) - 1; i >= 0; i-- {
			b.WriteString(fmt.Sprintf(" %s %s\n", dimRedStyle.Sprint("→"), grey.Sprint(m.CallStack[i])))
		}
	}

	return b.String()
}

// String returns a formatted string representation of the Memento.
func (m *Memento) String() string {
	return m.Error()
}
