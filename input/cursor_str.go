// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/neogeny/ogopego/util"
	"github.com/neogeny/ogopego/util/pyre"
)

// CursorHeavy holds the configuration and shared resources for a cursor.
type CursorHeavy struct {
	mu           sync.Mutex
	IgnoreCase   bool
	NameGuard    bool
	NameChars    string
	Source       string
	Patterns     TokenizingPatterns
	PatternCache map[string]pyre.Pattern
}

// StrCursor is a string-based implementation of the Cursor interface.
type StrCursor struct {
	text   string
	offset int
	heavy  *CursorHeavy
}

// NewStrCursor creates a new StrCursor with default configuration.
func NewStrCursor(text string) *StrCursor {
	return &StrCursor{
		text:   text,
		offset: 0,
		heavy: &CursorHeavy{
			NameGuard: false,
			Source:    "some input",
			Patterns:  DefaultPatterns(),
		},
	}
}

// NewStrCursorFromSource creates a new StrCursor starting at a specific offset.
func NewStrCursorFromSource(source, text string, start int) *StrCursor {
	if start > len(text) {
		start = len(text)
	}
	for start < len(text) && !utf8.RuneStart(text[start]) {
		start++
	}
	return &StrCursor{
		text:   text,
		offset: start,
		heavy: &CursorHeavy{
			NameGuard: false,
			Source:    source,
			Patterns:  DefaultPatterns(),
		},
	}
}

// Configure updates the cursor configuration.
func (s *StrCursor) Configure(cfg Cfg) {
	if cfg.Source != "" {
		s.heavy.Source = cfg.Source
	}

	s.heavy.mu.Lock()
	s.heavy.IgnoreCase = cfg.IgnoreCase
	s.heavy.NameChars = cfg.NameChars
	s.heavy.Patterns.Configure(cfg)

	if cfg.NameGuard != nil {
		s.heavy.NameGuard = *cfg.NameGuard ||
			s.heavy.NameChars != ""
	} else {
		s.heavy.NameGuard = cfg.NameChars != "" ||
			s.heavy.Patterns.NonDefault &&
				s.heavy.Patterns.Wsp != nil &&
				!s.heavy.Patterns.Wsp.IsEmpty()
	}
	s.heavy.mu.Unlock()
}

// InputSource returns the name or description of the input source.
func (s *StrCursor) InputSource() string {
	return s.heavy.Source
}

// Len returns the total input size in marks.
func (s *StrCursor) Len() int {
	return len(s.text)
}

// LineCount returns the total number of lines in the input.
func (s *StrCursor) LineCount() int {
	count := 0
	for range strings.Lines(s.text) {
		count++
	}
	return count
}

// LineAt returns the content of line n (0-indexed), or "" if out of range.
func (s *StrCursor) LineAt(n int) string {
	i := 0
	for line := range strings.Lines(s.text) {
		if i == n {
			return line
		}
		i++
	}
	return ""
}

// LinesAt returns lines [start, end) (0-indexed, half-open), or nil if empty.
func (s *StrCursor) LinesAt(start, end int) []string {
	if end <= start || start < 0 {
		return nil
	}
	out := make([]string, 0, end-start)
	i := 0
	for line := range strings.Lines(s.text) {
		if i >= end {
			break
		}
		if i >= start {
			out = append(out, line)
		}
		i++
	}
	return out
}

// Mark returns the current byte offset in the input.
func (s *StrCursor) Mark() int {
	return s.offset
}

// Reset sets the current byte offset in the input.
func (s *StrCursor) Reset(mark int) {
	s.offset = mark
}

// AsStr returns the entire input text.
func (s *StrCursor) AsStr() string {
	return s.text
}

// AsRef returns a reference to the input (the text itself for string cursors).
func (s *StrCursor) AsRef() string {
	return s.text
}

// IgnoreCase returns true if the cursor should ignore case during matching.
func (s *StrCursor) IgnoreCase() bool {
	return s.heavy.IgnoreCase
}

// NameGuard returns true if name guarding is enabled.
func (s *StrCursor) NameGuard() bool {
	return s.heavy.NameGuard
}

func (s *StrCursor) Lookahead(start int) string {
	if start >= len(s.text) {
		return ""
	}
	tail := s.text[start:]
	for line := range strings.Lines(tail) {
		return strings.TrimRight(line, "\n\r\t")
	}
	return tail
}

// AtEnd returns true if the cursor has reached the end of the input.
func (s *StrCursor) AtEnd() bool {
	return s.offset >= len(s.text)
}

// Next consumes and returns the next rune in the input.
func (s *StrCursor) Next() (rune, bool) {
	r, ok := s.Peek()
	if ok {
		s.offset += utf8.RuneLen(r)
	}
	return r, ok
}

// Peek returns the next rune in the input without consuming it.
func (s *StrCursor) Peek() (rune, bool) {
	if s.AtEnd() {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(s.text[s.offset:])
	return r, true
}

// PeekToken checks if the given token matches at the current offset without consuming it.
func (s *StrCursor) PeekToken(token string) bool {
	if s.offset+len(token) > len(s.text) {
		return false
	}
	slice := s.text[s.offset : s.offset+len(token)]
	if s.IgnoreCase() {
		return strings.EqualFold(slice, token)
	}
	return slice == token
}

// IsNameChar returns true if the rune can be part of a name.
func (s *StrCursor) IsNameChar(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c) ||
		strings.ContainsRune(s.heavy.NameChars, c)
}

// IsName returns true if the given string is a valid name.
func (s *StrCursor) IsName(token string) bool {
	if token == "" {
		return false
	}
	runes := []rune(token)
	first := runes[0]
	if first != '_' && !unicode.IsLetter(first) &&
		!strings.ContainsRune(s.heavy.NameChars, first) {
		return false
	}
	for _, r := range runes[1:] {
		if !s.IsNameChar(r) {
			return false
		}
	}
	return true
}

// MatchToken consumes the given token if it matches at the current offset.
// If nameguard is enabled and the token is a name, it checks the next character
// to ensure it is not also a name character (prevents partial name matches).
func (s *StrCursor) MatchToken(token string) bool {
	if !s.PeekToken(token) {
		return false
	}
	mark := s.offset
	s.offset += len(token)
	if s.heavy.NameGuard && s.IsName(token) {
		if r, size := utf8.DecodeRuneInString(s.text[s.offset:]); size > 0 && s.IsNameChar(r) {
			s.offset = mark
			return false
		}
	}
	return true
}

// MatchPattern matches a regular expression at the current offset and consumes it.
func (s *StrCursor) MatchPattern(pattern string) (string, bool) {
	pat := s.GetPattern(pattern)
	if pat == nil {
		return "", false
	}
	text := s.text[s.offset:]
	m, ok := pat.Match(text)
	if !ok {
		return "", false
	}
	s.offset += m.End()
	if g, ok := m.Group(1); ok {
		return g, true
	}
	if g, ok := m.Group(0); ok {
		return g, true
	}
	return "", false
}

// GetPattern compiles and caches a regular expression pattern.
func (s *StrCursor) GetPattern(pattern string) pyre.Pattern {
	s.heavy.mu.Lock()
	defer s.heavy.mu.Unlock()
	if s.heavy.PatternCache == nil {
		s.heavy.PatternCache = make(map[string]pyre.Pattern)
	}
	if p, ok := s.heavy.PatternCache[pattern]; ok {
		return p
	}
	p, err := pyre.Compile(pattern)
	if err != nil {
		return nil
	}
	s.heavy.PatternCache[pattern] = p
	return p
}

// MatchEOL matches an end-of-line (including trailing whitespace) and consumes it.
func (s *StrCursor) MatchEOL() bool {
	mark := s.offset
	s.eatSpacesNoNewlines()
	if n, ok := takeLinebreakLen(s.text[s.offset:]); ok {
		s.offset += n
		s.eatSpacesNoNewlines()
		return true
	}
	s.offset = mark
	return false
}

// NextToken consumes whitespace, comments, and end-of-line markers.
func (s *StrCursor) NextToken() {
	wsp := s.heavy.Patterns.Wsp
	eol := s.heavy.Patterns.Eol
	cmt := s.heavy.Patterns.Cmt

	for {
		prev := s.offset
		s.eatPattern(wsp)
		if s.eatPattern(eol) {
			s.eatPattern(wsp)
		}
		s.eatPattern(cmt)
		if s.AtEnd() || s.offset == prev {
			break
		}
	}
}

// Pos provides the "editor position" and the current byte offset into the input.
func (s *StrCursor) Pos() (int, int) {
	return s.PosAt(s.offset)
}

// PosAt provides the "editor position" and the given byte offset into the input.
func (s *StrCursor) PosAt(mark int) (int, int) {
	if mark <= 0 || len(s.text) == 0 {
		return 1, 1
	}
	if mark > len(s.text) {
		mark = len(s.text)
	}

	var line string
	lineno := 0 // NOTE: only empty strings render no Lines()
	for l := range strings.Lines(s.text[0:mark]) {
		line = l
		lineno += 1
	}
	line = util.ExpandTabs(line)
	col := len(util.StripRight(line))
	if col < len(line) {
		// mark was at the end of the line
		lineno += 1
		col = 1
	} else if col <= 0 {
		// mark was at the beginning of the line
		col = 1
	}
	return lineno, col
}

// Location returns the full location (source, line, col) at the current offset.
func (s *StrCursor) Location() Location {
	return s.LocationAt(s.offset)
}

// LocationAt returns the full location at the given offset.
func (s *StrCursor) LocationAt(mark int) Location {
	line, col := s.PosAt(mark)
	return Location{
		Source: s.InputSource(),
		Line:   line,
		Col:    col,
	}
}

// SetIgnoreCase updates the ignore-case setting.
func (s *StrCursor) SetIgnoreCase(ignore bool) {
	s.heavy.IgnoreCase = ignore
}

// SetPatterns updates the tokenizing patterns.
func (s *StrCursor) SetPatterns(patterns *TokenizingPatterns) {
	if patterns == nil {
		s.heavy.Patterns = DefaultPatterns()
	} else {
		s.heavy.Patterns = *patterns
	}
}

// Clone creates a copy of the cursor at the current offset.
func (s *StrCursor) Clone() Cursor {
	return &StrCursor{
		text:   s.text,
		offset: s.offset,
		heavy:  s.heavy,
	}
}

func (s *StrCursor) eatPattern(pat pyre.Pattern) bool {
	if pat == nil || s.AtEnd() || pat.Pattern() == "" {
		return false
	}
	text := s.text[s.offset:]
	if m, ok := pat.Match(text); ok && m.End() > 0 {
		s.offset += m.End()
		return true
	}
	return false
}

func (s *StrCursor) eatSpacesNoNewlines() {
	for {
		prev := s.offset
		s.offset += takeNonNewlineWhitespaceLen(s.text[s.offset:])
		if s.eatPattern(s.heavy.Patterns.Eol) {
			s.offset += takeNonNewlineWhitespaceLen(s.text[s.offset:])
		}
		s.eatPattern(s.heavy.Patterns.Cmt)
		if s.offset == prev {
			break
		}
	}
}

func takeLinebreakLen(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	switch s[0] {
	case '\n':
		return 1, true
	case '\r':
		if len(s) > 1 && s[1] == '\n' {
			return 2, true
		}
		return 1, true
	}
	return 0, false
}

func takeNonNewlineWhitespaceLen(s string) int {
	for i, c := range s {
		if c != ' ' && c != '\t' && c != '\f' {
			return i
		}
	}
	return len(s)
}
