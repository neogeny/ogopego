package context

import (
	"github.com/neogeny/ogopego/trees"
)

// Ctx represents the parsing context used during parse operations. It
// abstracts cursor position, memoization, failure reporting and tracing.
// See package documentation for details on life-cycle and semantics.
//
// The interface is implemented by CoreCtx returned from NewCtx.
type Ctx interface {
	Configurable

	Clone() Ctx
	Merge(other Ctx)
	// Cursor returns the input cursor.
	Cursor() Cursor
	// CallStack returns the current call stack.
	CallStack() CallStack
	// Tracer returns the tracer for the context.
	Tracer() Tracer
	// Mark returns the current position in the input.
	Mark() int
	// Reset resets the input position to the given mark.
	Reset(mark int)
	// AtEnd checks if the cursor is at the end of the input.
	AtEnd() bool
	// Next advances the cursor and returns the next rune.
	Next() (rune, bool)
	// Peek returns the next rune without advancing the cursor.
	Peek() (rune, bool)
	// Dot matches any character.
	Dot() (rune, error)
	// NextToken advances the cursor to the next token.
	NextToken()
	// MatchEOL matches an end-of-line.
	MatchEOL() bool
	// MatchToken matches a specific token string.
	MatchToken(token string) bool
	// MatchPattern matches a regular expression pattern.
	MatchPattern(pattern string) (string, error)
	Void()
	// EnterLookahead increments the lookahead depth counter.
	EnterLookahead()
	// LeaveLookahead decrements the lookahead depth counter.
	LeaveLookahead()
	// Fail records a parsing failure.
	Fail() error
	// Eof checks if the end of file has been reached.
	Eof() bool
	// EofCheck checks for end of file and returns an error if not at EOF.
	EofCheck() error
	// EolCheck checks for end of line and returns an error if not at EOL.
	EolCheck() error
	// Constant creates a tree node for a constant literal.
	Constant(literal any) (trees.Tree, error)
	// Enter pushes a rule onto the call stack for tracing.
	Enter(name string)
	// Leave pops a rule from the call stack.
	Leave()
	// Failure creates a new Nope error.
	Failure(start int, source error) error
	// FurthestFailure returns the furthest failure encountered so far.
	FurthestFailure() *ParseFailure
	// SetFurthestFailure sets the furthest failure.
	SetFurthestFailure(dis *ParseFailure)
	// IsKeyword checks if a name is a reserved keyword.
	IsKeyword(name string) bool
	// Intern interns a string to save memory.
	Intern(s string) string
	// ParseEOF checks if the parser should expect EOF.
	ParseEOF() bool
	// HeartbeatTick sends a tick to the heartbeat.
	HeartbeatTick()
	// Key creates a memoization key.
	Key(name string, canMemo bool) MemoKey
	// Memo retrieves a memoized result.
	Memo(key MemoKey) (Memo, bool)
	// Memoize stores a result in the memoization table.
	Memoize(key MemoKey, tree trees.Tree, mark int)
	// TrackRecursionDepth tracks recursion depth for a given key.
	TrackRecursionDepth(key MemoKey) error
	// Untrack removes a key from recursion tracking.
	Untrack(key MemoKey)

	// Cut marks a "cut" point in the parsing process.
	Cut()
	// IsCutSeen checks if a cut has been seen.
	IsCutSeen() bool
	// CutStackPush pushes the current cut state onto the stack.
	CutStackPush()
	// CutStackPop pops the cut state from the stack. Returns the popped value.
	CutStackPop() bool

	// ApplySemantics apply the semantics of Tree transformations in the Cfg for the grammar
	ApplySemantics(node trees.Tree, ruleName string, params []string) (trees.Tree, bool)
}
