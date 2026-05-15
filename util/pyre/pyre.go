// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package pyre

func Compile(pattern string) (Pattern, error) {
	return NewPCRE2CgoPattern(pattern)
}

func MustCompile(pattern string) Pattern {
	p, err := NewPCRE2CgoPattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
}
