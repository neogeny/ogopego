// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package pyre

import (
	"strings"

	"github.com/dlclark/regexp2"
)

type Regexp2Pattern struct {
	re *regexp2.Regexp
}

type Regexp2Match struct {
	match *regexp2.Match
	text  string
	re    *regexp2.Regexp
}

func NewRegexp2Pattern(pattern string) (*Regexp2Pattern, error) {
	re, err := regexp2.Compile(pattern, 0)
	if err != nil {
		return nil, err
	}
	return &Regexp2Pattern{re: re}, nil
}

func (p *Regexp2Pattern) Match(text string) (Match, bool) {
	m, err := p.re.FindStringMatchStartingAt(text, 0)
	if err != nil || m == nil {
		return nil, false
	}
	if rtoByte(text, m.Index) != 0 {
		return nil, false
	}
	return &Regexp2Match{match: m, text: text, re: p.re}, true
}

func (p *Regexp2Pattern) Search(text string) (Match, bool) {
	m, err := p.re.FindStringMatch(text)
	if err != nil || m == nil {
		return nil, false
	}
	return &Regexp2Match{match: m, text: text, re: p.re}, true
}

func (p *Regexp2Pattern) FullMatch(text string) (Match, bool) {
	m, err := p.re.FindStringMatchStartingAt(text, 0)
	if err != nil || m == nil {
		return nil, false
	}
	start := rtoByte(text, m.Index)
	end := rtoByte(text, m.Index+m.Length)
	if start != 0 || end != len(text) {
		return nil, false
	}
	return &Regexp2Match{match: m, text: text, re: p.re}, true
}

func (p *Regexp2Pattern) Split(text string, maxSplit int) []string {
	var result []string
	lastEnd := 0
	splitsDone := 0

	m, err := p.re.FindStringMatch(text)
	for m != nil && err == nil {
		if maxSplit > 0 && splitsDone >= maxSplit {
			break
		}
		start := rtoByte(text, m.Index)
		result = append(result, text[lastEnd:start])
		for i := 1; i < m.GroupCount(); i++ {
			g := m.GroupByNumber(i)
			if g.Index >= 0 {
				gs := rtoByte(text, g.Index)
				ge := rtoByte(text, g.Index+g.Length)
				result = append(result, text[gs:ge])
			} else {
				result = append(result, "")
			}
		}
		lastEnd = rtoByte(text, m.Index+m.Length)
		splitsDone++
		m, err = p.re.FindNextMatch(m)
	}

	result = append(result, text[lastEnd:])
	return result
}

func (p *Regexp2Pattern) FindAll(text string) [][]string {
	var result [][]string

	m, err := p.re.FindStringMatch(text)
	for m != nil && err == nil {
		result = append(result, extractGroups(m, text))
		m, err = p.re.FindNextMatch(m)
	}

	return result
}

func extractGroups(m *regexp2.Match, text string) []string {
	numGroups := m.GroupCount()
	if numGroups == 1 {
		return []string{m.String()}
	}
	if numGroups == 2 {
		g := m.GroupByNumber(1)
		if g.Index >= 0 {
			return []string{g.String()}
		}
		return []string{""}
	}
	groups := make([]string, numGroups-1)
	for i := 1; i < numGroups; i++ {
		g := m.GroupByNumber(i)
		if g.Index >= 0 {
			groups[i-1] = g.String()
		}
	}
	return groups
}

func (p *Regexp2Pattern) FindIter(text string) []Match {
	var result []Match

	m, err := p.re.FindStringMatch(text)
	for m != nil && err == nil {
		result = append(result, &Regexp2Match{match: m, text: text, re: p.re})
		m, err = p.re.FindNextMatch(m)
	}

	return result
}

func (p *Regexp2Pattern) Sub(repl, text string, count int) string {
	s, _ := p.SubN(repl, text, count)
	return s
}

func (p *Regexp2Pattern) SubN(repl, text string, count int) (string, int) {
	var result strings.Builder
	lastEnd := 0
	replacements := 0

	m, err := p.re.FindStringMatch(text)
	for m != nil && err == nil {
		if count > 0 && replacements >= count {
			break
		}
		start := rtoByte(text, m.Index)
		result.WriteString(text[lastEnd:start])
		result.WriteString(repl)
		lastEnd = rtoByte(text, m.Index+m.Length)
		replacements++
		m, err = p.re.FindNextMatch(m)
	}

	result.WriteString(text[lastEnd:])
	return result.String(), replacements
}

func (p *Regexp2Pattern) Pattern() string {
	return p.re.String()
}

func (p *Regexp2Pattern) MatchesEmpty() bool {
	m, err := p.re.MatchString("")
	return err == nil && m
}

func (p *Regexp2Pattern) IsEmpty() bool {
	return p == nil || strings.TrimSpace(p.re.String()) == ""
}

func (p *Regexp2Pattern) GroupIndex() map[string]int {
	names := p.re.GetGroupNames()
	result := make(map[string]int, len(names))
	for i, name := range names {
		if name != "" {
			result[name] = i
		}
	}
	return result
}

func (p *Regexp2Pattern) GroupsCount() int {
	// Find a match to get the group count; use empty string as a quick check
	m, _ := p.re.FindStringMatch("")
	if m != nil {
		return m.GroupCount() - 1
	}
	// If the pattern doesn't match "", try a simple "a"
	m, _ = p.re.FindStringMatch("a")
	if m != nil {
		return m.GroupCount() - 1
	}
	return 0
}

func (m *Regexp2Match) Group(i int) (string, bool) {
	g := m.match.GroupByNumber(i)
	if g == nil || g.Index < 0 {
		return "", false
	}
	return g.String(), true
}

func (m *Regexp2Match) Groups() []*string {
	groups := make([]*string, m.match.GroupCount()-1)
	for i := 1; i < m.match.GroupCount(); i++ {
		g := m.match.GroupByNumber(i)
		if g.Index >= 0 {
			groups[i-1] = new(g.String())
		}
	}
	return groups
}

func (m *Regexp2Match) Start() int {
	return rtoByte(m.text, m.match.Index)
}

func (m *Regexp2Match) End() int {
	return rtoByte(m.text, m.match.Index+m.match.Length)
}

func (m *Regexp2Match) Span() (int, int) {
	return m.Start(), m.End()
}

func (m *Regexp2Match) GroupName(name string) (string, bool) {
	g := m.match.GroupByName(name)
	if g == nil || g.Index < 0 {
		return "", false
	}
	return g.String(), true
}

func (m *Regexp2Match) GroupDict() map[string]*string {
	groups := m.match.Groups()
	result := make(map[string]*string)
	for _, g := range groups {
		if g.Name != "" {
			if g.Index >= 0 {
				result[g.Name] = new(g.String())
			}
		}
	}
	return result
}

func (m *Regexp2Match) Expand(template string) string {
	result, _ := m.re.Replace(m.match.String(), template, 0, 1)
	return result
}

func rtoByte(s string, runePos int) int {
	i := 0
	for pos := range s {
		if i == runePos {
			return pos
		}
		i++
	}
	return len(s)
}
