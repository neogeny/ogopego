// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type prettyWriter struct {
	buf    strings.Builder
	indent int
	amount int
}

func newPrettyWriter() *prettyWriter {
	return &prettyWriter{amount: 4}
}

func (w *prettyWriter) WriteLine(s string) {
	if s == "" {
		w.buf.WriteByte('\n')
		return
	}
	pad := strings.Repeat(" ", w.indent*w.amount)
	for _, line := range strings.Split(s, "\n") {
		w.buf.WriteString(pad)
		w.buf.WriteString(line)
		w.buf.WriteByte('\n')
	}
}

func (w *prettyWriter) Indent() { w.indent++ }
func (w *prettyWriter) Dedent() { w.indent-- }
func (w *prettyWriter) Reset()  { w.buf.Reset() }
func (w *prettyWriter) String() string {
	return strings.TrimRight(w.buf.String(), "\n")
}

const pep8llen = 72

// Token

func (m *Token) PrettyPrint() string {
	return fmt.Sprintf(`"%s"`, m.Token)
}

// Pattern

func (m *Pattern) PrettyPrint() string {
	if strings.Contains(m.Pattern, "/") {
		return fmt.Sprintf(`?"%s"`, m.Pattern)
	}
	return fmt.Sprintf("/%s/", m.Pattern)
}

// Constant

func (m *Constant) PrettyPrint() string {
	if strings.Count(m.Literal, "\n") <= 1 {
		return fmt.Sprintf("`%s`", m.Literal)
	}
	return fmt.Sprintf("```%s```", m.Literal)
}

// Alert

func (m *Alert) PrettyPrint() string {
	return fmt.Sprintf("%s`%s`", strings.Repeat("^", m.Level), m.Literal)
}

// Call

func (m *Call) PrettyPrint() string {
	return m.Name
}

// RuleInclude

func (m *RuleInclude) PrettyPrint() string {
	return fmt.Sprintf(">%s", m.Name)
}

// Leaf terminals

func (m *Cut) PrettyPrint() string          { return "~" }
func (m *Dot) PrettyPrint() string          { return "." }
func (m *EOF) PrettyPrint() string          { return "$" }
func (m *EOL) PrettyPrint() string          { return "$->" }
func (m *Fail) PrettyPrint() string         { return "!()" }
func (m *NULL) PrettyPrint() string         { return "" }
func (m *Void) PrettyPrint() string         { return "()" }
func (m *EmptyClosure) PrettyPrint() string { return "{}" }

// Option

func (m *Option) PrettyPrint() string {
	return m.Exp.PrettyPrint()
}

// Group

func (m *Group) PrettyPrint() string {
	inner := m.Exp.PrettyPrint()
	if strings.ContainsRune(inner, '\n') {
		w := newPrettyWriter()
		w.WriteLine("(")
		w.Indent()
		w.WriteLine(inner)
		w.Dedent()
		w.WriteLine(")")
		return w.String()
	}
	return fmt.Sprintf("(%s)", inner)
}

// SkipGroup

func (m *SkipGroup) PrettyPrint() string {
	return fmt.Sprintf("(?:%s)", m.Exp.PrettyPrint())
}

// Lookahead

func (m *Lookahead) PrettyPrint() string {
	return fmt.Sprintf("&%s", m.Exp.PrettyPrint())
}

// NegativeLookahead

func (m *NegativeLookahead) PrettyPrint() string {
	return fmt.Sprintf("!%s", m.Exp.PrettyPrint())
}

// SkipTo

func (m *SkipTo) PrettyPrint() string {
	return fmt.Sprintf("->%s", m.Exp.PrettyPrint())
}

// Optional

func (m *Optional) PrettyPrint() string {
	return fmt.Sprintf("[%s]", m.Exp.PrettyPrint())
}

// Closure

func (m *Closure) PrettyPrint() string {
	return fmt.Sprintf("{%s}*", m.Exp.PrettyPrint())
}

// PositiveClosure

func (m *PositiveClosure) PrettyPrint() string {
	return fmt.Sprintf("{%s}+", m.Exp.PrettyPrint())
}

// Override

func (m *Override) PrettyPrint() string {
	return fmt.Sprintf("=%s", m.Exp.PrettyPrint())
}

// OverrideList

func (m *OverrideList) PrettyPrint() string {
	return fmt.Sprintf("+=%s", m.Exp.PrettyPrint())
}

// Synth

func (m *Synth) PrettyPrint() string {
	return m.Exp.PrettyPrint()
}

// Named

func (m *Named) PrettyPrint() string {
	return fmt.Sprintf("%s=%s", m.Name, m.Exp.PrettyPrint())
}

// NamedList

func (m *NamedList) PrettyPrint() string {
	return fmt.Sprintf("%s+=%s", m.Name, m.Exp.PrettyPrint())
}

// Join

func (m *Join) PrettyPrint() string {
	return fmt.Sprintf("%s%%{%s}*", m.Sep.PrettyPrint(), m.Exp.PrettyPrint())
}

// PositiveJoin

func (m *PositiveJoin) PrettyPrint() string {
	return fmt.Sprintf("%s%%{%s}+", m.Sep.PrettyPrint(), m.Exp.PrettyPrint())
}

// Gather

func (m *Gather) PrettyPrint() string {
	return fmt.Sprintf("%s.{%s}*", m.Sep.PrettyPrint(), m.Exp.PrettyPrint())
}

// PositiveGather

func (m *PositiveGather) PrettyPrint() string {
	return fmt.Sprintf("%s.{%s}+", m.Sep.PrettyPrint(), m.Exp.PrettyPrint())
}

// Sequence

func (m *Sequence) PrettyPrint() string {
	items := make([]string, 0, len(m.Sequence))
	for _, exp := range m.Sequence {
		if _, ok := exp.(*EOF); ok {
			continue
		}
		items = append(items, exp.PrettyPrint())
	}
	single := strings.Join(items, " ")
	hasMulti := false
	for _, s := range items {
		if strings.ContainsRune(s, '\n') {
			hasMulti = true
			break
		}
	}
	if !hasMulti && len(single) <= pep8llen {
		return single
	}
	w := newPrettyWriter()
	for _, item := range items {
		w.WriteLine(item)
	}
	return w.String()
}

// Choice

func (m *Choice) PrettyPrint() string {
	opts := make([]string, len(m.Options))
	for i, opt := range m.Options {
		opts[i] = opt.Exp.PrettyPrint()
	}
	hasMulti := false
	for _, s := range opts {
		if strings.ContainsRune(s, '\n') {
			hasMulti = true
			break
		}
	}
	singleLine := strings.Join(opts, " | ")
	if !hasMulti && len(singleLine) <= pep8llen {
		return singleLine
	}
	w := newPrettyWriter()
	for _, opt := range opts {
		w.WriteLine("| " + opt)
	}
	return w.String()
}

// Rule

func (m *Rule) PrettyPrint() string {
	w := newPrettyWriter()
	if m.NoStak {
		w.WriteLine("@nostak")
	}
	if m.NoMemo {
		w.WriteLine("@nomemo")
	}
	if m.IsName {
		w.WriteLine("@name")
	}
	params := ""
	if len(m.Params) > 0 {
		params = fmt.Sprintf("[%s]", strings.Join(m.Params, ", "))
	}
	exp := m.Exp.PrettyPrint()
	if strings.ContainsRune(exp, '\n') {
		w.WriteLine(fmt.Sprintf("%s%s:", m.Name, params))
		w.Indent()
		w.WriteLine(exp)
		w.Dedent()
	} else {
		w.WriteLine(fmt.Sprintf("%s%s: %s", m.Name, params, exp))
	}
	return w.String()
}

// Grammar

func directiveValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "True"
		}
		return "False"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (m *Grammar) PrettyPrint() string {
	w := newPrettyWriter()
	w.WriteLine(fmt.Sprintf("@@grammar :: %s", m.Name))

	if m.Directives == nil {
		m.Directives = orderedmap.New()
	}

	knownDir := func(k string) (any, bool) { return m.Directives.Get(k) }

	if dir, ok := knownDir("whitespace"); ok {
		w.WriteLine(fmt.Sprintf("@@whitespace :: /%s/", directiveValue(dir)))
	}
	if dir, ok := knownDir("comments"); ok {
		w.WriteLine(fmt.Sprintf("@@comments :: /%s/", directiveValue(dir)))
	}
	if dir, ok := knownDir("eol_comments"); ok {
		w.WriteLine(fmt.Sprintf("@@eol_comments :: /%s/", directiveValue(dir)))
	}
	if dir, ok := knownDir("namechars"); ok {
		w.WriteLine(fmt.Sprintf("@@namechars :: \"%s\"", directiveValue(dir)))
	}
	if dir, ok := knownDir("ignorecase"); ok {
		if v, ok := dir.(bool); ok && v {
			w.WriteLine("@@ignorecase :: True")
		}
	}
	if dir, ok := knownDir("nameguard"); ok {
		if v, ok := dir.(bool); ok && v {
			w.WriteLine("@@nameguard :: True")
		}
	}
	if dir, ok := knownDir("left_recursion"); ok {
		if v, ok := dir.(bool); ok && !v {
			w.WriteLine("@@left_recursion :: False")
		}
	}
	if dir, ok := knownDir("parseinfo"); ok {
		if v, ok := dir.(bool); ok && !v {
			w.WriteLine("@@parseinfo :: False")
		}
	}
	if dir, ok := knownDir("memoization"); ok {
		if v, ok := dir.(bool); ok && !v {
			w.WriteLine("@@memoization :: False")
		}
	}

	known := map[string]bool{
		"grammar": true, "whitespace": true, "comments": true, "eol_comments": true,
		"namechars": true, "ignorecase": true, "nameguard": true,
		"left_recursion": true, "parseinfo": true, "memoization": true,
	}
	for _, k := range m.Directives.Keys() {
		if known[k] {
			continue
		}
		v, _ := m.Directives.Get(k)
		w.WriteLine(fmt.Sprintf("@@%s :: %s", k, directiveValue(v)))
	}

	if len(m.Keywords) > 0 {
		w.WriteLine("")
		for _, chunk := range chunkStrings(m.Keywords, 8) {
			w.WriteLine(fmt.Sprintf("@@keyword :: %s", strings.Join(chunk, " ")))
		}
	}

	for _, rule := range m.Rules {
		w.WriteLine("")
		w.WriteLine(rule.PrettyPrint())
	}
	return w.String()
}

func chunkStrings(s []string, size int) [][]string {
	if len(s) == 0 {
		return nil
	}
	var chunks [][]string
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
