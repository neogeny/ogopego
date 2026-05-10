package input

type Location struct {
	Source string
	Line   int
	Col    int
}

type MatchResult interface {
	End() int
	Group(i int) (string, bool)
}

type Pattern interface {
	Match(text string) (MatchResult, bool)
	Pattern() string
}

type TokenizingPatterns struct {
	Wsp Pattern
	Cmt Pattern
	Eol Pattern
}

type Cursor interface {
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
	MatchPattern(pattern Pattern) (string, bool)
	MatchEOL() bool
	NextToken()
	Pos() (int, int)
	PosAt(mark int) (int, int)
	Location() Location
	LocationAt(mark int) Location
	SetPatterns(patterns *TokenizingPatterns)
	Clone() Cursor
}
