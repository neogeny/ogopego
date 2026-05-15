// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package input

import (
	"strings"
	"unicode/utf8"

	"github.com/neogeny/ogopego/util/pyre"
)

type CursorHeavy struct {
	IgnoreCase bool
	NameGuard  bool
	Source     string
	Patterns   *TokenizingPatterns
}

type StrCursor struct {
	text   string
	offset int
	heavy  *CursorHeavy
}

func NewStrCursor(text string) *StrCursor {
	return &StrCursor{
		text:   text,
		offset: 0,
		heavy: &CursorHeavy{
			NameGuard: true,
			Source:    "some input",
			Patterns:  &TokenizingPatterns{},
		},
	}
}

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
			NameGuard: true,
			Source:    source,
			Patterns:  &TokenizingPatterns{},
		},
	}
}

func (s *StrCursor) Configure(cfg Cfg) {
	s.heavy.IgnoreCase = cfg.IgnoreCase
	s.heavy.NameGuard = cfg.NameGuard
	if cfg.Source != "" {
		s.heavy.Source = cfg.Source
	}

	if cfg.Whitespace != nil {
		if *cfg.Whitespace != "" {
			if pat, err := pyre.Compile(*cfg.Whitespace); err == nil {
				s.heavy.Patterns.Wsp = pat
			}
		} else {
			s.heavy.Patterns.Wsp = nil
		}
	}
	if cfg.Comments != "" {
		if pat, err := pyre.Compile(cfg.Comments); err == nil {
			s.heavy.Patterns.Cmt = pat
		}
	}
	if cfg.EolComments != "" {
		if pat, err := pyre.Compile(cfg.EolComments); err == nil {
			s.heavy.Patterns.Eol = pat
		}
	}
}

func (s *StrCursor) InputSource() string {
	return s.heavy.Source
}

func (s *StrCursor) Mark() int {
	return s.offset
}

func (s *StrCursor) Reset(mark int) {
	s.offset = mark
}

func (s *StrCursor) AsStr() string {
	return s.text
}

func (s *StrCursor) AsRef() string {
	return s.text
}

func (s *StrCursor) IgnoreCase() bool {
	return s.heavy.IgnoreCase
}

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

func (s *StrCursor) AtEnd() bool {
	return s.offset >= len(s.text)
}

func (s *StrCursor) Next() (rune, bool) {
	r, ok := s.Peek()
	if ok {
		s.offset += utf8.RuneLen(r)
	}
	return r, ok
}

func (s *StrCursor) Peek() (rune, bool) {
	if s.AtEnd() {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(s.text[s.offset:])
	return r, true
}

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

func (s *StrCursor) MatchToken(token string) bool {
	if s.PeekToken(token) {
		s.offset += len(token)
		return true
	}
	return false
}

func (s *StrCursor) MatchPattern(pat pyre.Pattern) (string, bool) {
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

func (s *StrCursor) NextToken() {
	if s.heavy.Patterns == nil {
		return
	}

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

func (s *StrCursor) Pos() (int, int) {
	return s.PosAt(s.offset)
}

func tabDisplayWidth(s string) int {
	var w int
	for _, r := range s {
		if r == '\t' {
			w += 4
		} else {
			w++
		}
	}
	return w
}

func (s *StrCursor) PosAt(mark int) (int, int) {
	if mark > len(s.text) {
		mark = len(s.text)
	}
	if mark <= 0 {
		return 0, 0
	}
	lineno := 0
	var line string
	for l := range strings.Lines(s.text[0:mark]) {
		line = l
		lineno += 1
	}
	return lineno, tabDisplayWidth(line)
}

func (s *StrCursor) Location() Location {
	return s.LocationAt(s.offset)
}

func (s *StrCursor) LocationAt(mark int) Location {
	line, col := s.PosAt(mark)
	return Location{
		Source: s.InputSource(),
		Line:   line,
		Col:    col,
	}
}

func (s *StrCursor) SetIgnoreCase(ignore bool) {
	s.heavy = &CursorHeavy{
		IgnoreCase: ignore,
		NameGuard:  s.heavy.NameGuard,
		Source:     s.heavy.Source,
		Patterns:   s.heavy.Patterns,
	}
}

func (s *StrCursor) SetPatterns(patterns *TokenizingPatterns) {
	s.heavy = &CursorHeavy{
		IgnoreCase: s.heavy.IgnoreCase,
		NameGuard:  s.heavy.NameGuard,
		Source:     s.heavy.Source,
		Patterns:   patterns,
	}
}

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
