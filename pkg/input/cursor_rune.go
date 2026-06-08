package input

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/neogeny/ogopego/pkg/util"
	pyre2 "github.com/neogeny/ogopego/pkg/util/pyre"
)

// RuneCursor is a Cursor implementation that stores pre-computed runes for
// efficient rune-level operations (Next, Peek, etc.).
//
// Marks are rune indices, not byte offsets. This avoids the per-call
// []rune conversion overhead in regexp matching and provides the
// foundation for a future rune-based regexp2 API.
type RuneCursor struct {
	runes   []rune
	runePos int
	heavy   *CursorHeavy
}

// NewRuneCursor creates a new RuneCursor for the given text.
func NewRuneCursor(text string) *RuneCursor {
	return newRuneCursor(text, "", 0)
}

// NewRuneCursorFromSource creates a new RuneCursor with a source name and start offset.
func NewRuneCursorFromSource(source, text string, start int) *RuneCursor {
	return newRuneCursor(text, source, start)
}

func newRuneCursor(text, source string, start int) *RuneCursor {
	if start > len(text) {
		start = len(text)
	}
	for start < len(text) && !utf8.RuneStart(text[start]) {
		start++
	}
	runes := []rune(text)
	rp := 0
	bytePos := 0
	for rp < len(runes) && bytePos < start {
		bytePos += utf8.RuneLen(runes[rp])
		rp++
	}
	return &RuneCursor{
		runes:   runes,
		runePos: rp,
		heavy: &CursorHeavy{
			NameGuard: false,
			Source:    source,
			Patterns:  DefaultPatterns(),
		},
	}
}

// --- Configurable ---

func (c *RuneCursor) Configure(cfg Cfg) {
	if cfg.Source != "" {
		c.heavy.Source = cfg.Source
	}
	c.heavy.mu.Lock()
	c.heavy.IgnoreCase = cfg.IgnoreCase
	c.heavy.NameChars = cfg.NameChars
	c.heavy.Patterns.Configure(cfg)
	if cfg.NameGuard != nil {
		c.heavy.NameGuard = *cfg.NameGuard || c.heavy.NameChars != ""
	} else {
		c.heavy.NameGuard = cfg.NameChars != "" ||
			c.heavy.Patterns.NonDefault &&
				c.heavy.Patterns.Wsp != nil &&
				!c.heavy.Patterns.Wsp.IsEmpty()
	}
	c.heavy.mu.Unlock()
}

// --- Cursor interface ---

func (c *RuneCursor) InputSource() string {
	return c.heavy.Source
}

func (c *RuneCursor) Mark() int {
	return c.runePos
}

func (c *RuneCursor) Reset(mark int) {
	if mark < 0 {
		mark = 0
	}
	if mark > len(c.runes) {
		mark = len(c.runes)
	}
	c.runePos = mark
}

func (c *RuneCursor) Len() int {
	return len(c.runes)
}

func (c *RuneCursor) LineCount() int {
	count := 0
	for _, r := range c.runes {
		if r == '\n' {
			count++
		}
	}
	if len(c.runes) > 0 && c.runes[len(c.runes)-1] != '\n' {
		count++
	}
	return count
}

func (c *RuneCursor) LineAt(n int) string {
	lineno := 0
	start := 0
	for i := 0; i < len(c.runes); i++ {
		if c.runes[i] == '\n' {
			if lineno == n {
				return string(c.runes[start : i+1])
			}
			lineno++
			start = i + 1
		}
	}
	if lineno == n {
		return string(c.runes[start:])
	}
	return ""
}

func (c *RuneCursor) LinesAt(start, end int) []string {
	if end <= start || start < 0 {
		return nil
	}
	out := make([]string, 0, end-start)
	lineno := 0
	lineStart := 0
	for i := 0; i < len(c.runes) && lineno < end; i++ {
		if c.runes[i] == '\n' {
			if lineno >= start {
				out = append(out, string(c.runes[lineStart:i+1]))
			}
			lineno++
			lineStart = i + 1
		}
	}
	if lineno >= start && lineno < end {
		if len(c.runes) > 0 && c.runes[len(c.runes)-1] != '\n' {
			out = append(out, string(c.runes[lineStart:]))
		}
	}
	return out
}

func (c *RuneCursor) AsStr() string {
	return string(c.runes)
}

func (c *RuneCursor) AsRef() string {
	return string(c.runes)
}

func (c *RuneCursor) IgnoreCase() bool {
	return c.heavy.IgnoreCase
}

func (c *RuneCursor) NameGuard() bool {
	return c.heavy.NameGuard
}

func (c *RuneCursor) Lookahead(start int) string {
	if start >= len(c.runes) {
		return ""
	}
	end := start
	for end < len(c.runes) && c.runes[end] != '\n' {
		end++
	}
	line := string(c.runes[start:end])
	return strings.TrimRight(line, "\n\r\t")
}

func (c *RuneCursor) AtEnd() bool {
	return c.runePos >= len(c.runes)
}

func (c *RuneCursor) Next() (rune, bool) {
	if c.AtEnd() {
		return 0, false
	}
	r := c.runes[c.runePos]
	c.runePos++
	return r, true
}

func (c *RuneCursor) Peek() (rune, bool) {
	if c.AtEnd() {
		return 0, false
	}
	return c.runes[c.runePos], true
}

func (c *RuneCursor) PeekToken(token string) bool {
	if c.IgnoreCase() {
		i := 0
		for _, r := range token {
			if c.runePos+i >= len(c.runes) {
				return false
			}
			if unicode.ToLower(c.runes[c.runePos+i]) != unicode.ToLower(r) {
				return false
			}
			i++
		}
		return true
	}
	i := 0
	for _, r := range token {
		if c.runePos+i >= len(c.runes) {
			return false
		}
		if c.runes[c.runePos+i] != r {
			return false
		}
		i++
	}
	return true
}

func (c *RuneCursor) IsNameChar(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) ||
		strings.ContainsRune(c.heavy.NameChars, r)
}

func (c *RuneCursor) IsName(token string) bool {
	if token == "" {
		return false
	}
	first := true
	for _, r := range token {
		if first {
			if r != '_' && !unicode.IsLetter(r) &&
				!strings.ContainsRune(c.heavy.NameChars, r) {
				return false
			}
			first = false
		} else if !c.IsNameChar(r) {
			return false
		}
	}
	return true
}

func (c *RuneCursor) MatchToken(token string) bool {
	if !c.PeekToken(token) {
		return false
	}
	mark := c.runePos
	c.runePos += utf8.RuneCountInString(token)
	if c.heavy.NameGuard && c.IsName(token) {
		if c.runePos < len(c.runes) && c.IsNameChar(c.runes[c.runePos]) {
			c.runePos = mark
			return false
		}
	}
	return true
}

func (c *RuneCursor) MatchPattern(pattern string) (string, bool) {
	pat := c.GetPattern(pattern)
	if pat == nil {
		return "", false
	}
	m, ok := pat.MatchRunes(c.runes[c.runePos:])
	if !ok {
		return "", false
	}
	c.runePos += m.End()
	if g, ok := m.Group(1); ok {
		return g, true
	}
	if g, ok := m.Group(0); ok {
		return g, true
	}
	return "", false
}

func (c *RuneCursor) MatchName() (string, bool) {
	if c.AtEnd() {
		return "", false
	}
	mark := c.runePos
	first := c.runes[mark]
	if first != '_' && !unicode.IsLetter(first) &&
		!strings.ContainsRune(c.heavy.NameChars, first) {
		return "", false
	}
	c.runePos++
	for c.runePos < len(c.runes) && c.IsNameChar(c.runes[c.runePos]) {
		c.runePos++
	}
	return string(c.runes[mark:c.runePos]), true
}

func (c *RuneCursor) MatchInt() (int, bool) {
	mark := c.runePos
	if !c.consumeSignedInt() {
		return 0, false
	}
	raw := string(c.runes[mark:c.runePos])
	n, err := strconv.Atoi(cleanNumber(raw))
	if err != nil {
		c.runePos = mark
		return 0, false
	}
	return n, true
}

func (c *RuneCursor) MatchUInt() (uint64, bool) {
	mark := c.runePos
	if !c.consumeUInt() {
		return 0, false
	}
	raw := string(c.runes[mark:c.runePos])
	n, err := strconv.ParseUint(cleanNumber(raw), 10, 64)
	if err != nil {
		c.runePos = mark
		return 0, false
	}
	return n, true
}

func (c *RuneCursor) consumeSignedInt() bool {
	mark := c.runePos
	c.consumeSign()
	if !c.consumeUInt() {
		c.runePos = mark
		return false
	}
	return true
}

func (c *RuneCursor) consumeSign() bool {
	if c.runePos < len(c.runes) && (c.runes[c.runePos] == '+' || c.runes[c.runePos] == '-') {
		c.runePos++
		return true
	}
	return false
}

func (c *RuneCursor) consumeUInt() bool {
	start := c.runePos
	for c.runePos < len(c.runes) {
		r := c.runes[c.runePos]
		if r >= '0' && r <= '9' {
			c.runePos++
		} else if r == '_' {
			if c.runePos+1 < len(c.runes) && c.runes[c.runePos+1] >= '0' && c.runes[c.runePos+1] <= '9' {
				c.runePos++
			} else {
				c.runePos = start
				return false
			}
		} else if unicode.IsLetter(r) {
			c.runePos = start
			return false
		} else {
			break
		}
	}
	return c.runePos != start
}

func (c *RuneCursor) MatchFloat() (float64, bool) {
	mark := c.runePos
	if !c.consumeSignedInt() {
		c.runePos = mark
		return 0, false
	}
	if c.runePos < len(c.runes) && c.runes[c.runePos] == '.' {
		c.runePos++
		c.consumeUInt()
	}
	if c.runePos < len(c.runes) && (c.runes[c.runePos] == 'e' || c.runes[c.runePos] == 'E') {
		expMark := c.runePos
		c.runePos++
		if !c.consumeSignedInt() {
			c.runePos = expMark
		}
	}
	raw := string(c.runes[mark:c.runePos])
	f, err := strconv.ParseFloat(cleanNumber(raw), 64)
	if err != nil {
		c.runePos = mark
		return 0, false
	}
	return f, true
}

func (c *RuneCursor) MatchBool() (bool, bool) {
	if c.AtEnd() {
		return false, false
	}
	mark := c.runePos
	first := unicode.ToLower(c.runes[mark])
	if first == 't' {
		if mark+4 <= len(c.runes) && string(c.runes[mark+1:mark+4]) == "rue" {
			c.runePos = mark + 4
			return true, true
		}
	} else if first == 'f' {
		if mark+5 <= len(c.runes) && string(c.runes[mark+1:mark+5]) == "alse" {
			c.runePos = mark + 5
			return false, true
		}
	}
	return false, false
}

func (c *RuneCursor) GetPattern(pattern string) pyre2.Pattern {
	c.heavy.mu.Lock()
	defer c.heavy.mu.Unlock()
	if c.heavy.PatternCache == nil {
		c.heavy.PatternCache = make(map[string]pyre2.Pattern)
	}
	if p, ok := c.heavy.PatternCache[pattern]; ok {
		return p
	}
	p, err := pyre2.Compile(pattern)
	if err != nil {
		return nil
	}
	c.heavy.PatternCache[pattern] = p
	return p
}

func (c *RuneCursor) MatchEOL() bool {
	mark := c.runePos
	c.eatSpacesNoNewlines()
	if n, ok := takeLinebreakLenRunes(c.runes, c.runePos); ok {
		c.runePos += n
		c.eatSpacesNoNewlines()
		return true
	}
	c.runePos = mark
	return false
}

func (c *RuneCursor) NextToken() {
	wsp := c.heavy.Patterns.Wsp
	eol := c.heavy.Patterns.Eol
	cmt := c.heavy.Patterns.Cmt
	for {
		prev := c.runePos
		c.eatPattern(wsp)
		if c.eatPattern(eol) {
			c.eatPattern(wsp)
		}
		c.eatPattern(cmt)
		if c.AtEnd() || c.runePos == prev {
			break
		}
	}
}

func (c *RuneCursor) eatPattern(pat pyre2.Pattern) bool {
	if pat == nil || c.AtEnd() || pat.Pattern() == "" {
		return false
	}
	m, ok := pat.MatchRunes(c.runes[c.runePos:])
	if ok && m.End() > 0 {
		c.runePos += m.End()
		return true
	}
	return false
}

func (c *RuneCursor) eatSpacesNoNewlines() {
	for {
		prev := c.runePos
		c.skipNonNewlineWhitespace()
		if c.eatPattern(c.heavy.Patterns.Eol) {
			c.skipNonNewlineWhitespace()
		}
		c.eatPattern(c.heavy.Patterns.Cmt)
		if c.runePos == prev {
			break
		}
	}
}

func (c *RuneCursor) skipNonNewlineWhitespace() {
	for c.runePos < len(c.runes) {
		r := c.runes[c.runePos]
		if r != ' ' && r != '\t' && r != '\f' {
			break
		}
		c.runePos++
	}
}

func takeLinebreakLenRunes(runes []rune, pos int) (int, bool) {
	if pos >= len(runes) {
		return 0, false
	}
	switch runes[pos] {
	case '\n':
		return 1, true
	case '\r':
		if pos+1 < len(runes) && runes[pos+1] == '\n' {
			return 2, true
		}
		return 1, true
	}
	return 0, false
}

func (c *RuneCursor) Pos() (int, int) {
	return c.PosAt(c.runePos)
}

func (c *RuneCursor) PosAt(mark int) (int, int) {
	if mark <= 0 || len(c.runes) == 0 {
		return 1, 1
	}
	if mark > len(c.runes) {
		mark = len(c.runes)
	}
	// Delegate to the same algorithm as StrCursor by converting the prefix to string.
	// PosAt is only called for error reporting, not on hot paths.
	return strPosAt(string(c.runes[:mark]))
}

func strPosAt(prefix string) (int, int) {
	var line string
	lineno := 0
	for l := range strings.Lines(prefix) {
		line = l
		lineno++
	}
	line = util.ExpandTabs(line)
	col := len(util.StripRight(line))
	if col < len(line) {
		lineno++
		col = 1
	} else if col <= 0 {
		col = 1
	}
	return lineno, col
}

func (c *RuneCursor) Location() Location {
	line, col := c.Pos()
	return Location{Source: c.heavy.Source, Line: line, Col: col}
}

func (c *RuneCursor) LocationAt(mark int) Location {
	line, col := c.PosAt(mark)
	return Location{Source: c.heavy.Source, Line: line, Col: col}
}

func (c *RuneCursor) SetPatterns(patterns *TokenizingPatterns) {
	if patterns == nil {
		c.heavy.Patterns = DefaultPatterns()
	} else {
		c.heavy.Patterns = *patterns
	}
}

func (c *RuneCursor) SetIgnoreCase(ignore bool) {
	c.heavy.IgnoreCase = ignore
}

func (c *RuneCursor) Clone() Cursor {
	return &RuneCursor{
		runes:   c.runes,
		runePos: c.runePos,
		heavy:   c.heavy,
	}
}

var _ Cursor = (*RuneCursor)(nil)

func cleanNumber(raw string) string { return strings.ReplaceAll(raw, "_", "") }
