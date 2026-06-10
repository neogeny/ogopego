// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego"
)

var bootGrammar *Grammar

// BootGrammar returns the internal boot grammar that is used to parse EBNF
// grammar strings. The result is cached after the first call.
func BootGrammar() (*Grammar, error) {
	if bootGrammar == nil {
		boot, err := loadBootGrammar()
		if err != nil {
			return nil, err
		}
		bootGrammar = boot
	}
	return bootGrammar, nil
}

// loadBootGrammar parses and initializes a boot grammar from JSON bytes.
func loadBootGrammar() (*Grammar, error) {
	g, err := LoadGrammarFromJSON(ogopego.TatSuGrammarJSON)
	if err != nil {
		return nil, err
	}
	if err := g.Initialize(); err != nil {
		return nil, err
	}
	// these semantics apply only to parsers that parse grammars
	g.Semantics = EBNFGrammarSemantics{}
	return g, nil
}
