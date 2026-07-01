// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0
package context

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/util"
)

var (
	errStyle  = color.New(color.FgRed, color.Bold)
	blueStyle = color.New(color.FgBlue, color.Bold)
	bold      = color.New(color.Bold)
	grey      = color.New(color.FgHiBlack)
	dimStyle  = color.New(color.Faint, color.FgWhite)
)

// Memento captures the state of the parser at a specific point for error reporting.
type Memento struct {
	Cursor    input.Cursor
	Msg       string
	Start     int
	Mark      int
	CallStack CallStack
}

// NewMemento constructs a Memento capturing cursor state and a message for
// later diagnostic reporting.
func NewMemento(start int, msg string, cursor input.Cursor, callstack CallStack) Memento {
	return Memento{
		Cursor:    cursor,
		Msg:       msg,
		Start:     start,
		Mark:      cursor.Mark(),
		CallStack: callstack,
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
	bluePipe := blueStyle.Sprint("│")
	arrow := blueStyle.Sprint("─→")

	errMsg := fmt.Sprintf("%s: %s\n", errLabel, bold.Sprint(m.Msg))
	b.WriteString("\n")
	b.WriteString(errMsg)
	b.WriteString(fmt.Sprintf("  %s %s @ [%d:%d]\n", arrow, m.Cursor.InputSource(), line, col))

	b.WriteString(fmt.Sprintf(" %5s%s\n", "", bluePipe))
	start := line - 4
	if start < 0 {
		start = 0
	}
	for i, linestr := range m.Cursor.LinesAt(start, line+1) {
		linestr = util.StripRight(linestr)
		disp := util.ExpandTabs(linestr)
		lineno := dimStyle.Sprintf("%5d", start+i)
		b.WriteString(fmt.Sprintf("%s %s %s\n", lineno, bluePipe, disp))
	}
	pad := strings.Repeat(" ", col)
	_, _ = fmt.Fprintf(&b, " %5s%s %s\n", "", bluePipe, errStyle.Sprintf("%s⌃ %s", pad, errMsg))

	if !m.CallStack.IsEmpty() {
		b.WriteString("\n")
		for name := range m.CallStack.All() {
			b.WriteString(fmt.Sprintf(" %s %s\n", errStyle.Sprint("→"), grey.Sprint(name)))
		}
	}

	return b.String()
}

// String returns a formatted string representation of the Memento.
func (m *Memento) String() string {
	return m.Error()
}
