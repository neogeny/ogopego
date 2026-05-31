// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package pyre

import (
	"regexp"
	"strings"
)

type GoPattern struct {
	re *regexp.Regexp
}

type GoMatch struct {
	text  string
	match []int
	re    *regexp.Regexp
}

func NewGoPattern(pattern string) (*GoPattern, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &GoPattern{re: re}, nil
}

func (p *GoPattern) Match(text string) (Match, bool) {
	loc := p.re.FindStringSubmatchIndex(text)
	if loc == nil || loc[0] != 0 {
		return nil, false
	}
	return &GoMatch{text: text, match: loc, re: p.re}, true
}

func (p *GoPattern) MatchRunes(runes []rune) (Match, bool) {
	return p.Match(string(runes))
}

func (p *GoPattern) Search(text string) (Match, bool) {
	loc := p.re.FindStringSubmatchIndex(text)
	if loc == nil {
		return nil, false
	}
	return &GoMatch{text: text, match: loc, re: p.re}, true
}

func (p *GoPattern) FullMatch(text string) (Match, bool) {
	loc := p.re.FindStringSubmatchIndex(text)
	if loc == nil {
		return nil, false
	}
	if loc[0] != 0 || loc[1] != len(text) {
		return nil, false
	}
	return &GoMatch{text: text, match: loc, re: p.re}, true
}

func (p *GoPattern) Split(text string, maxSplit int) []string {
	var result []string
	lastEnd := 0
	splitsDone := 0

	for lastEnd <= len(text) {
		if maxSplit > 0 && splitsDone >= maxSplit {
			break
		}
		loc := p.re.FindStringSubmatchIndex(text[lastEnd:])
		if loc == nil {
			break
		}
		for i := range loc {
			if loc[i] >= 0 {
				loc[i] += lastEnd
			}
		}
		result = append(result, text[lastEnd:loc[0]])
		for i := 2; i < len(loc); i += 2 {
			if loc[i] >= 0 {
				result = append(result, text[loc[i]:loc[i+1]])
			} else {
				result = append(result, "")
			}
		}
		lastEnd = loc[1]
		if loc[0] == loc[1] {
			lastEnd++
		}
		splitsDone++
	}

	result = append(result, text[lastEnd:])
	return result
}

func (p *GoPattern) FindAll(text string) [][]string {
	matches := p.re.FindAllStringSubmatch(text, -1)
	if matches == nil {
		return nil
	}
	result := make([][]string, len(matches))
	if p.re.NumSubexp() == 0 {
		for i, m := range matches {
			result[i] = []string{m[0]}
		}
	} else {
		for i, m := range matches {
			result[i] = m[1:]
		}
	}
	return result
}

func (p *GoPattern) FindIter(text string) []Match {
	locs := p.re.FindAllStringSubmatchIndex(text, -1)
	if locs == nil {
		return nil
	}
	result := make([]Match, len(locs))
	for i, loc := range locs {
		result[i] = &GoMatch{text: text, match: loc, re: p.re}
	}
	return result
}

func (p *GoPattern) Sub(repl, text string, count int) string {
	s, _ := p.SubN(repl, text, count)
	return s
}

func (p *GoPattern) SubN(repl, text string, count int) (string, int) {
	var result strings.Builder
	lastEnd := 0
	replacements := 0

	for {
		if count > 0 && replacements >= count {
			break
		}
		loc := p.re.FindStringSubmatchIndex(text[lastEnd:])
		if loc == nil {
			break
		}
		for i := range loc {
			if loc[i] >= 0 {
				loc[i] += lastEnd
			}
		}
		result.WriteString(text[lastEnd:loc[0]])
		result.WriteString(repl)
		lastEnd = loc[1]
		if loc[0] == loc[1] {
			lastEnd++
		}
		if lastEnd > len(text) {
			break
		}
		replacements++
	}

	result.WriteString(text[lastEnd:])
	return result.String(), replacements
}

func (p *GoPattern) Pattern() string {
	return p.re.String()
}

func (p *GoPattern) MatchesEmpty() bool {
	return p.re.MatchString("")
}

func (p *GoPattern) IsEmpty() bool {
	if p == nil || p.re == nil {
		return true
	}
	return strings.TrimSpace(p.re.String()) == ""
}

func (p *GoPattern) GroupIndex() map[string]int {
	names := p.re.SubexpNames()
	result := make(map[string]int, len(names))
	for i, name := range names {
		if name != "" {
			result[name] = i
		}
	}
	return result
}

func (p *GoPattern) GroupsCount() int {
	return p.re.NumSubexp()
}

func (m *GoMatch) Group(i int) (string, bool) {
	if 2*i+1 >= len(m.match) || m.match[2*i] < 0 {
		return "", false
	}
	return m.text[m.match[2*i]:m.match[2*i+1]], true
}

func (m *GoMatch) Groups() []*string {
	n := (len(m.match) / 2) - 1
	if n <= 0 {
		return nil
	}
	result := make([]*string, n)
	for i := 1; i <= n; i++ {
		if m.match[2*i] >= 0 {
			s := m.text[m.match[2*i]:m.match[2*i+1]]
			result[i-1] = &s
		}
	}
	return result
}

func (m *GoMatch) Start() int {
	return m.match[0]
}

func (m *GoMatch) End() int {
	return m.match[1]
}

func (m *GoMatch) Span() (int, int) {
	return m.match[0], m.match[1]
}

func (m *GoMatch) GroupName(name string) (string, bool) {
	idx := m.re.SubexpIndex(name)
	if idx < 0 {
		return "", false
	}
	return m.Group(idx)
}

func (m *GoMatch) GroupDict() map[string]*string {
	names := m.re.SubexpNames()
	result := make(map[string]*string)
	for i, name := range names {
		if name == "" {
			continue
		}
		if 2*i+1 < len(m.match) && m.match[2*i] >= 0 {
			s := m.text[m.match[2*i]:m.match[2*i+1]]
			result[name] = &s
		}
	}
	return result
}

func (m *GoMatch) Expand(template string) string {
	return string(m.re.ExpandString(nil, template, m.text, m.match))
}
