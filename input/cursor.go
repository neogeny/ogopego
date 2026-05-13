package input

import "github.com/neogeny/ogopego/util/pyre"

type Location struct {
	Source string
	Line   int
	Col    int
}

type TokenizingPatterns struct {
	Wsp pyre.Pattern
	Cmt pyre.Pattern
	Eol pyre.Pattern
}

type Cursor interface {
	Configurable
	InputSource() string
	Mark() int
	Reset(mark int)
	AsStr() string
	AsRef() string
	IgnoreCase() bool
	NameGuard() bool
	Lookahead(start int) string
	AtEnd() bool
	Next() (rune, bool)
	Peek() (rune, bool)
	PeekToken(token string) bool
	MatchToken(token string) bool
	MatchPattern(pattern pyre.Pattern) (string, bool)
	MatchEOL() bool
	NextToken()
	Pos() (int, int)
	PosAt(mark int) (int, int)
	Location() Location
	LocationAt(mark int) Location
	SetPatterns(patterns *TokenizingPatterns)
	SetIgnoreCase(ignore bool)
	Clone() Cursor
}
