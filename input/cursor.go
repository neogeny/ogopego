// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import "github.com/neogeny/ogopego/util/pyre"

// Location represents a specific point in the input source.
type Location struct {
	Source string
	Line   int
	Col    int
}

// Cursor is an abstraction over input sources providing position, lookahead
// and tokenization utilities used by the parsing runtime.
type Cursor interface {
	Configurable
	// InputSource returns the name of the input source.
	InputSource() string
	// Mark returns the current position of the cursor.
	Mark() int
	// Reset sets the cursor position to the given mark.
	Reset(mark int)
	// Len returns the total input size in marks.
	Len() int
	// LineCount returns the total number of lines in the input.
	LineCount() int
	// LineAt returns the content of line n (0-indexed), or "" if out of range.
	LineAt(n int) string
	// LinesAt returns lines [start, end) (0-indexed, half-open), or nil if empty.
	LinesAt(start, end int) []string
	// AsStr returns the entire input as a string.
	AsStr() string
	// AsRef returns a reference to the input string.
	AsRef() string
	// IgnoreCase returns true if case should be ignored during matching.
	IgnoreCase() bool
	// NameGuard returns true if name guards should be enforced.
	NameGuard() bool
	// Lookahead returns a substring from the given start mark to the current cursor position.
	Lookahead(start int) string
	// AtEnd returns true if the cursor is at the end of the input.
	AtEnd() bool
	// Next advances the cursor and returns the next rune and true, or 0 and false if at end.
	Next() (rune, bool)
	// Peek returns the next rune without advancing the cursor.
	Peek() (rune, bool)
	// PeekToken checks if the next token matches the given string without advancing the cursor.
	PeekToken(token string) bool
	// MatchToken attempts to match the given token and advances the cursor if successful.
	MatchToken(token string) bool
	// IsNameChar returns true if the rune can be part of a name.
	IsNameChar(c rune) bool
	// IsName returns true if the given string is a valid name.
	IsName(s string) bool
	// MatchPattern attempts to match the given pattern and advances the cursor if successful.
	MatchPattern(pattern string) (string, bool)
	// GetPattern compiles and caches a regular expression pattern.
	GetPattern(pattern string) pyre.Pattern
	// MatchEOL attempts to match an end-of-line sequence and advances the cursor if successful.
	MatchEOL() bool
	// NextToken advances the cursor past the current token (whitespace and comments).
	NextToken()
	// Pos returns the current line and column number.
	Pos() (int, int)
	// PosAt returns the line and column number at a given mark.
	PosAt(mark int) (int, int)
	// Location returns the current source location.
	Location() Location
	// LocationAt returns the source location at a given mark.
	LocationAt(mark int) Location
	// SetPatterns sets the tokenizing patterns for the cursor.
	SetPatterns(patterns *TokenizingPatterns)
	// SetIgnoreCase sets whether the cursor should ignore case.
	SetIgnoreCase(ignore bool)
	// Clone creates a copy of the cursor.
	Clone() Cursor
}
