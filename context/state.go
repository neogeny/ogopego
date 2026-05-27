package context

import (
	"github.com/neogeny/ogopego/input"
)

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

// Cursor returns the input cursor associated with this ParseState.
func (ps *ParseState) Cursor() input.Cursor {
	return ps.cursor
}

// CutSeen returns true if a cut has been encountered in this ParseState.
func (ps *ParseState) CutSeen() bool {
	return ps.cutSeen
}

// SetCut marks that a cut has been encountered in this ParseState.
func (ps *ParseState) SetCut() {
	ps.cutSeen = true
}

// LastNode returns the last parsed node in this ParseState.
func (ps *ParseState) LastNode() any {
	return ps.lastNode
}

// SetLastNode sets the last parsed node for this ParseState.
func (ps *ParseState) SetLastNode(n any) {
	ps.lastNode = n
}

// Clone creates a deep copy of the ParseState.
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

// Top returns the current (topmost) ParseState on the stack.
func (ss *StateStack) Top() *ParseState {
	return ss.states[len(ss.states)-1]
}

// Push creates a new ParseState by cloning the current top and pushes it onto the stack.
func (ss *StateStack) Push() *ParseState {
	s := ss.Top().Clone()
	ss.states = append(ss.states, s)
	return s
}

// Pop removes and returns the topmost ParseState from the stack.
func (ss *StateStack) Pop() *ParseState {
	popped := ss.states[len(ss.states)-1]
	ss.states = ss.states[:len(ss.states)-1]
	return popped
}

// Undo is an alias for Pop, removing the topmost ParseState.
func (ss *StateStack) Undo() *ParseState {
	return ss.Pop()
}

// Merge merges the topmost ParseState into the one below it, effectively discarding the top.
func (ss *StateStack) Merge() {
	if len(ss.states) < 2 {
		return
	}
	child := ss.Pop()
	ss.Top().cursor.Reset(child.cursor.Mark())
	ss.Top().cutSeen = child.cutSeen
}
