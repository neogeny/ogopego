// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package pyre

// LookaheadSupport reports whether the compiled regex engine supports
// lookahead assertions ((?=...), (?!...)). Always false in this configuration.
const LookaheadSupport = false

func Compile(pattern string) (Pattern, error) {
	return NewRegexp2V2Pattern(pattern)
}

func MustCompile(pattern string) Pattern {
	p, err := NewRegexp2V2Pattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
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
