// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import "github.com/neogeny/ogopego"

var bootGrammar *Grammar

func BootGrammar() (*Grammar, error) {
	if bootGrammar == nil {
		boot, err := LoadBootGrammar()
		if err != nil {
			return nil, err
		}
		bootGrammar = boot
	}
	return bootGrammar, nil
}

// LoadBootGrammar parses and initializes a boot grammar from JSON bytes.
func LoadBootGrammar() (*Grammar, error) {
	g, err := ParseGrammar(ogopego.TatSuGrammarJSON)
	if err != nil {
		return nil, err
	}
	if err := g.Initialize(); err != nil {
		return nil, err
	}
	return g, nil
}
