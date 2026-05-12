package context

import (
	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

type Ctx interface {
	Cursor() input.Cursor
	CallStack() CallStack
	Tracer() Tracer
	Mark() int
	Reset(mark int)
	AtEnd() bool
	Next() (rune, bool)
	Peek() (rune, bool)
	Dot() (rune, error)
	NextToken()
	MatchEOL() bool
	MatchToken(token string) bool
	MatchPattern(pattern string) (string, bool)
	GetPattern(pattern string) pyre.Pattern
	Token(token string) (string, error)
	Pattern(pattern string) (string, error)
	Void() error
	Fail() error
	Eof() bool
	EofCheck() error
	EolCheck() error
	Constant(literal any) (trees.Tree, error)
	Enter(name string)
	Leave()
	Failure(start int, source error) *DisasterReport
	FurthestFailure() *DisasterReport
	SetFurthestFailure(dis *DisasterReport)
	IsKeyword(name string) bool
	SetKeywords(keywords []string)
	Intern(s string) string
	ParseEOF() bool
	HeartbeatTick()
}
