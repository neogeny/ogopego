package util

import (
	"math"
	"strings"
)

const spaces = " \t\n\r\f"

func StripLeft(s string) string {
	return strings.TrimLeft(s, spaces)
}

func StripRight(s string) string {
	return strings.TrimRight(s, spaces)
}

func ExpandTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

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
