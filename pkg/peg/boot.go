// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego"
	"github.com/neogeny/ogopego/pkg/trees"
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

// EBNFGrammarSemantics are the library's semantics for parsing grammars
// in the syntax defined by tatsu.ebnf and tatsu.json.
type EBNFGrammarSemantics struct{}

func (EBNFGrammarSemantics) Apply(
	node trees.Tree,
	ruleName string,
	params []string) (trees.Tree, bool) {
	switch ruleName {
	case "true":
		return trees.TRUE, true
	case "false":
		return trees.FALSE, true
	case "null":
		return trees.NULL, true
	case "meta":
		text := textValue(node)
		switch text {
		case "name":
			return &trees.Node{TypeName: "NameMeta"}, true
		case "int":
			return &trees.Node{TypeName: "IntMeta"}, true
		case "uint":
			return &trees.Node{TypeName: "UIntMeta"}, true
		case "float":
			return &trees.Node{TypeName: "FloatMeta"}, true
		case "bool":
			return &trees.Node{TypeName: "BoolMeta"}, true
		}
	}
	// do nothing by default
	// return false to tell caller to apply default processing
	return node, false
}
