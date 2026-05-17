package util

import (
	"math"
	"strings"
)

const spaces = " \t\n\r\f"

// StripLeft removes leading whitespace characters from s.
func StripLeft(s string) string {
	return strings.TrimLeft(s, spaces)
}

// StripRight removes trailing whitespace characters from s.
func StripRight(s string) string {
	return strings.TrimRight(s, spaces)
}

// ExpandTabs replaces tabs with 4 spaces.
func ExpandTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

// Dedent strips common identation from a text block starting after the first line.
func Dedent(s string) string {
	indent := math.MaxInt
	skip := true
	for line := range strings.Lines(s) {
		if skip {
			skip = false
			continue
		}
		if len(StripLeft(line)) == 0 {
			continue
		}
		indent = min(indent, len(line)-len(StripLeft(line)))
	}

	if indent <= 0 || indent >= math.MaxInt {
		return s
	}

	var b strings.Builder
	first := true
	for line := range strings.Lines(s) {
		if first {
			first = false
			b.WriteString(line)
			continue
		}
		if len(StripLeft(line)) == 0 {
			b.WriteString(line)
		} else {
			b.WriteString(line[indent:])
		}
	}
	return b.String()
}

func TabDisplayWidth(s string) int {
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
