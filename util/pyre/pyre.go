// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package pyre

func Compile(pattern string) (Pattern, error) {
	return NewRegexp2Pattern(pattern)
}

func MustCompile(pattern string) Pattern {
	p, err := NewRegexp2Pattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
}
