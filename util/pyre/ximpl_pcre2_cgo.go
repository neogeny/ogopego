package pyre

import (
	"strings"

	"github.com/Jemmic/go-pcre2"
)

type PCRE2CgoPattern struct {
	re      *pcre2.Regexp
	pattern string
}

type PCRE2CgoMatch struct {
	m   *pcre2.Matcher
	re  *pcre2.Regexp
	idx []int
}

func NewPCRE2CgoPattern(pattern string) (*PCRE2CgoPattern, error) {
	re, err := pcre2.Compile(pattern, 0)
	if err != nil {
		return nil, err
	}
	return &PCRE2CgoPattern{re: re, pattern: pattern}, nil
}

func (p *PCRE2CgoPattern) Match(text string) (Match, bool) {
	m := p.re.MatcherString(text, 0)
	if !m.Matches() {
		m.Free()
		return nil, false
	}
	idx := m.Index()
	if idx[0] != 0 {
		m.Free()
		return nil, false
	}
	return &PCRE2CgoMatch{m: m, re: p.re, idx: idx}, true
}

func (p *PCRE2CgoPattern) Search(text string) (Match, bool) {
	m := p.re.MatcherString(text, 0)
	if !m.Matches() {
		m.Free()
		return nil, false
	}
	return &PCRE2CgoMatch{m: m, re: p.re, idx: m.Index()}, true
}

func (p *PCRE2CgoPattern) FullMatch(text string) (Match, bool) {
	m := p.re.MatcherString(text, 0)
	if !m.Matches() {
		m.Free()
		return nil, false
	}
	idx := m.Index()
	if idx[0] != 0 || idx[1] != len(text) {
		m.Free()
		return nil, false
	}
	return &PCRE2CgoMatch{m: m, re: p.re, idx: idx}, true
}

func (p *PCRE2CgoPattern) Split(text string, maxSplit int) []string {
	var result []string
	lastEnd := 0
	splitsDone := 0

	for lastEnd <= len(text) {
		if maxSplit > 0 && splitsDone >= maxSplit {
			break
		}
		m := p.re.MatcherString(text[lastEnd:], 0)
		if !m.Matches() {
			m.Free()
			break
		}
		idx := m.Index()
		start := lastEnd + idx[0]
		end := lastEnd + idx[1]

		result = append(result, text[lastEnd:start])
		ng := m.Groups()
		for i := 1; i <= ng; i++ {
			if m.Present(i) {
				result = append(result, m.GroupString(i))
			} else {
				result = append(result, "")
			}
		}
		m.Free()

		if start == end {
			lastEnd = end + 1
		} else {
			lastEnd = end
		}
		splitsDone++
		if lastEnd >= len(text) {
			break
		}
	}

	result = append(result, text[lastEnd:])
	return result
}

func (p *PCRE2CgoPattern) FindAll(text string) [][]string {
	var result [][]string
	pos := 0

	for pos <= len(text) {
		m := p.re.MatcherString(text[pos:], 0)
		if !m.Matches() {
			m.Free()
			break
		}
		idx := m.Index()
		ng := m.Groups()
		if ng == 0 {
			result = append(result, []string{m.GroupString(0)})
		} else {
			groups := make([]string, ng)
			for i := 1; i <= ng; i++ {
				if m.Present(i) {
					groups[i-1] = m.GroupString(i)
				}
			}
			result = append(result, groups)
		}
		m.Free()

		if idx[0] == idx[1] {
			pos++
		} else {
			pos += idx[1]
		}
		if pos > len(text) {
			break
		}
	}

	return result
}

func (p *PCRE2CgoPattern) FindIter(text string) []Match {
	var result []Match
	pos := 0

	for pos <= len(text) {
		m := p.re.MatcherString(text[pos:], 0)
		if !m.Matches() {
			m.Free()
			break
		}
		idx := m.Index()
		result = append(result, &PCRE2CgoMatch{m: m, re: p.re, idx: idx})

		if idx[0] == idx[1] {
			pos++
		} else {
			pos += idx[1]
		}
		if pos > len(text) {
			break
		}
	}

	return result
}

func (p *PCRE2CgoPattern) Sub(repl, text string, count int) string {
	s, _ := p.SubN(repl, text, count)
	return s
}

func (p *PCRE2CgoPattern) SubN(repl, text string, count int) (string, int) {
	var result strings.Builder
	lastEnd := 0
	replacements := 0

	for lastEnd <= len(text) {
		if count > 0 && replacements >= count {
			break
		}
		m := p.re.MatcherString(text[lastEnd:], 0)
		if !m.Matches() {
			m.Free()
			break
		}
		idx := m.Index()
		start := lastEnd + idx[0]
		end := lastEnd + idx[1]

		result.WriteString(text[lastEnd:start])
		result.WriteString(repl)

		lastEnd = end
		replacements++
		m.Free()

		if start == end {
			lastEnd++
		}
		if lastEnd > len(text) {
			break
		}
	}

	result.WriteString(text[lastEnd:])
	return result.String(), replacements
}

func (p *PCRE2CgoPattern) Pattern() string {
	return p.pattern
}

func (p *PCRE2CgoPattern) MatchesEmpty() bool {
	m := p.re.MatcherString("", 0)
	defer m.Free()
	return m.Matches()
}

func (p *PCRE2CgoPattern) IsEmpty() bool {
	return strings.TrimSpace(p.pattern) == ""
}

func (p *PCRE2CgoPattern) GroupIndex() map[string]int {
	return nil
}

func (p *PCRE2CgoPattern) GroupsCount() int {
	return p.re.Groups()
}

func (m *PCRE2CgoMatch) Group(i int) (string, bool) {
	if i < 0 || i > m.m.Groups() {
		return "", false
	}
	if m.m.Present(i) {
		return m.m.GroupString(i), true
	}
	return "", false
}

func (m *PCRE2CgoMatch) Groups() []*string {
	ng := m.m.Groups()
	if ng == 0 {
		return nil
	}
	result := make([]*string, ng)
	for i := 1; i <= ng; i++ {
		if m.m.Present(i) {
			s := m.m.GroupString(i)
			result[i-1] = &s
		}
	}
	return result
}

func (m *PCRE2CgoMatch) Start() int {
	return m.idx[0]
}

func (m *PCRE2CgoMatch) End() int {
	return m.idx[1]
}

func (m *PCRE2CgoMatch) Span() (int, int) {
	return m.idx[0], m.idx[1]
}

func (m *PCRE2CgoMatch) GroupName(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	s, err := m.m.NamedString(name)
	if err != nil {
		return "", false
	}
	return s, true
}

func (m *PCRE2CgoMatch) GroupDict() map[string]*string {
	return nil
}

func (m *PCRE2CgoMatch) Expand(template string) string {
	matchedText := m.m.GroupString(0)
	return m.re.ReplaceAllString(matchedText, template, 0)
}
