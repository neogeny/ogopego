// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// LoadBootGrammar parses and initializes a boot grammar from JSON bytes.
func LoadBootGrammar(data []byte) (*Grammar, error) {
	g, err := ParseGrammar(data)
	if err != nil {
		return nil, err
	}
	if err := g.Initialize(); err != nil {
		return nil, err
	}
	return g, nil
}
