package context

import (
	"github.com/neogeny/ogopego/input"
)

// CallStack is a slice of call-site names representing the parser call stack.
type CallStack []string

// ParseState holds lightweight parsing state for a cursor position.
type ParseState struct {
	cursor   input.Cursor
	cutSeen  bool
	lastNode any
}

// NewParseState creates a new ParseState for the provided cursor.
func NewParseState(cursor input.Cursor) *ParseState {
	return &ParseState{
		cursor: cursor.Clone(),
	}
}

func (ps *ParseState) Cursor() input.Cursor {
	return ps.cursor
}

func (ps *ParseState) CutSeen() bool {
	return ps.cutSeen
}

func (ps *ParseState) SetCut() {
	ps.cutSeen = true
}

func (ps *ParseState) LastNode() any {
	return ps.lastNode
}

func (ps *ParseState) SetLastNode(n any) {
	ps.lastNode = n
}

func (ps *ParseState) Clone() *ParseState {
	return &ParseState{
		cursor:   ps.cursor.Clone(),
		cutSeen:  ps.cutSeen,
		lastNode: ps.lastNode,
	}
}

// StateStack manages a stack of ParseState objects for nested parsing.
type StateStack struct {
	states []*ParseState
}

// NewStateStack returns a new StateStack initialized with a ParseState for
// the given cursor.
func NewStateStack(cursor input.Cursor) *StateStack {
	return &StateStack{
		states: []*ParseState{NewParseState(cursor)},
	}
}

func (ss *StateStack) Top() *ParseState {
	return ss.states[len(ss.states)-1]
}

func (ss *StateStack) Push() *ParseState {
	s := ss.Top().Clone()
	ss.states = append(ss.states, s)
	return s
}

func (ss *StateStack) Pop() *ParseState {
	popped := ss.states[len(ss.states)-1]
	ss.states = ss.states[:len(ss.states)-1]
	return popped
}

func (ss *StateStack) Undo() *ParseState {
	return ss.Pop()
}

func (ss *StateStack) Merge() {
	if len(ss.states) < 2 {
		return
	}
	child := ss.Pop()
	ss.Top().cursor.Reset(child.cursor.Mark())
	ss.Top().cutSeen = child.cutSeen
}
