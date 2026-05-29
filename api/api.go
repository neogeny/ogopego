// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

// Package api provides the public API for compiling grammars and parsing input.
//
// Most users will call Compile to create a Grammar from an EBNF string, then
// ParseInput to parse text with it. The cfg parameter controls tracing,
// colorization, whitespace handling, keywords, and other options.
//
// The JSON variant functions produce JSON-compatible output (map[string]any,
// []any, string, float64, bool, nil) suitable for json.Marshal.
package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/trees"
)

// Cfg is an alias for config.Cfg. It controls parsing behavior including
// tracing, colorization, whitespace handling, keywords, and more.
// Pass nil to use defaults.
type Cfg = config.Cfg

var (
	compileMu    sync.RWMutex
	compileCache = make(map[string]*peg.Grammar)
)

// BootGrammar returns the internal boot grammar that is used to parse EBNF
// grammar strings. The result is cached after the first call.
func BootGrammar() (*peg.Grammar, error) {
	return peg.BootGrammar()
}

// ParseGrammar parses a grammar string using the boot grammar and returns
// the raw parse tree. Use Compile instead to get a usable Grammar object.
func ParseGrammar(grammar string, cfg *Cfg) (trees.Tree, error) {
	grammar = strings.TrimRight(grammar, " \t\r\n")
	boot, err := BootGrammar()
	if err != nil {
		return nil, err
	}
	if boot.Semantics == nil {
		return nil, fmt.Errorf("boot grammar semantics not set")
	}

	cursor := input.NewStrCursor(grammar)

	directivesCfg := boot.CfgFromDirectives()
	if directivesCfg.Semantics == nil {
		// FIXME: this looks like debugging boot Grammar semantics
		return nil, fmt.Errorf("semantics not returned from grammar")
	}

	ctx := context.NewCtx(cursor, cfg)
	ctx.Configure(*directivesCfg)
	if ctx.Cfg().Semantics == nil {
		return nil, fmt.Errorf("boot semantics not passed to Ctx")
	}

	return boot.ParseAt(ctx, cfg)
}

// ParseGrammarToJSON parses a grammar string and returns the raw parse tree
// as a JSON-compatible value.
func ParseGrammarToJSON(grammar string, cfg *Cfg) (any, error) {
	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return nil, err
	}
	return asjson.AsJSON(tree), nil
}

// ParseGrammarToJSONStr is like ParseGrammarToJSON but returns a JSON
// string.
func ParseGrammarToJSONStr(grammar string, cfg *Cfg) (string, error) {
	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return "", err
	}
	return asjson.AsJSONStr(tree), nil
}

// Compile parses a grammar string and returns a compiled Grammar ready for
// parsing input. Results are cached by grammar string.
func Compile(grammar string, cfg *Cfg) (*peg.Grammar, error) {
	compileMu.RLock()
	if g, ok := compileCache[grammar]; ok {
		compileMu.RUnlock()
		return g, nil
	}
	compileMu.RUnlock()

	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return nil, err
	}
	g, err := peg.CompileGrammar(tree)
	if err != nil {
		return nil, err
	}

	compileMu.Lock()
	compileCache[grammar] = g
	compileMu.Unlock()
	return g, nil
}

// CompileToJSON parses a grammar string and returns the compiled Grammar as a
// JSON-compatible value.
func CompileToJSON(grammar string, cfg *Cfg) (any, error) {
	g, err := Compile(grammar, cfg)
	if err != nil {
		return nil, err
	}
	return asjson.AsJSON(g), nil
}

// CompileToJSONString is like CompileToJSON but returns a JSON string.
func CompileToJSONString(grammar string, cfg *Cfg) (string, error) {
	g, err := Compile(grammar, cfg)
	if err != nil {
		return "", err
	}
	return asjson.AsJSONStr(g), nil
}

// ParseInput parses the given text using a compiled Grammar and returns the
// resulting AST as a Tree value.
func ParseInput(parser *peg.Grammar, text string, cfg *Cfg) (trees.Tree, error) {
	ctx := context.NewCtx(input.NewStrCursor(text), cfg)
	return parser.ParseAt(ctx, cfg)
}

// ParseInputToJSON parses input text and returns the resulting AST as a
// JSON-compatible value.
func ParseInputToJSON(parser *peg.Grammar, text string, cfg *Cfg) (any, error) {
	tree, err := ParseInput(parser, text, cfg)
	if err != nil {
		return nil, err
	}
	return trees.TreeToJSON(tree), nil
}

// ParseInputToJSONString is like ParseInputToJSON but returns a JSON string.
func ParseInputToJSONString(parser *peg.Grammar, text string, cfg *Cfg) (string, error) {
	tree, err := ParseInput(parser, text, cfg)
	if err != nil {
		return "", err
	}
	return trees.TreeToJSONStr(tree), nil
}

// LoadGrammarFromJSON deserializes a Grammar from JSON output produced by
//
//	CompileToJSON or peg.serializeGrammar.
func LoadGrammarFromJSON(data []byte) (*peg.Grammar, error) {
	return peg.LoadGrammarFromJSON(data)
}
