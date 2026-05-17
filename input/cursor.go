// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package input

import "github.com/neogeny/ogopego/util/pyre"

type Location struct {
	Source string
	Line   int
	Col    int
}

// TokenizingPatterns groups precompiled patterns used by cursors for
// whitespace, comments and EOL detection.
type TokenizingPatterns struct {
	Wsp pyre.Pattern
	Cmt pyre.Pattern
	Eol pyre.Pattern
}

type Cursor interface {
	// Cursor is an abstraction over input sources providing position, lookahead
	// and tokenization utilities used by the parsing runtime.
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
