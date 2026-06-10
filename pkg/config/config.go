// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strings"

	"github.com/neogeny/ogopego/pkg/trees"
	"github.com/neogeny/ogopego/pkg/util"
	"github.com/neogeny/ogopego/pkg/util/heartbeat"
)

// ProgramName is the name of the CLI binary for use in messages.
//
//goland:noinspection GoUnusedConst
const ProgramName = "OGoPEGo"

// DefaultPerLineMemos is the default number of memo entries per input line.
const DefaultPerLineMemos = 8 // NOTE: Magic !!!

// Configurable is implemented by types that can be configured using a Cfg.
type Configurable interface {
	Configure(cfg Cfg)
}

type GrammarSemantics interface {
	Apply(node trees.Tree, ruleName string, params []string) (trees.Tree, bool)
}

// Cfg configures grammar compilation and input parsing. Use DefaultCfg() or
// pass nil to API functions for defaults. Individual fields override
// grammar-level @@directives.
type Cfg struct {
	Name   string // grammar name (overrides @@grammar)
	Source string // source description for error messages
	Start  string // start rule name (default: "start")

	Concurrency       bool    // enable concurrent Choice evaluation (experimental)
	NoMemo            bool    // disable memoization
	NoPruneMemosOnCut bool    // disable pruning memos on cut (~)
	PerLineMemos      float64 // memo entries per input line

	Trace    bool // enable parse trace output to stderr
	Colorize bool // colorize trace output

	Grammar         string // grammar text (overrides @@grammar body)
	NoLeftRecursion bool   // disable left-recursion support

	IgnoreCase bool   // case-insensitive matching
	NameChars  string // additional name characters (implies NameGuard)
	NameGuard  *bool  // enforce word boundaries on names

	Whitespace  *string // whitespace pattern (nil = default, &"" = none)
	Comments    string  // comment pattern
	EolComments string  // end-of-line comment pattern

	Keywords []string // reserved words

	ParseInfo bool // attach source position info to AST nodes

	Semantics GrammarSemantics
	Heart     heartbeat.Heart // progress callback (CLI progress bars)
}

// DefaultCfg returns the default configuration. Pass nil to API functions
// to use these defaults.
func DefaultCfg() *Cfg {
	return &Cfg{
		NoMemo:            false,
		NoPruneMemosOnCut: false,
		PerLineMemos:      DefaultPerLineMemos,
		Trace:             false,
		Colorize:          false,
		NoLeftRecursion:   false,
		IgnoreCase:        false,
		NameGuard:         nil,
		Whitespace:        nil,
		Keywords:          nil,
		ParseInfo:         false,
	}
}

// New returns a new Cfg produced by applying cfg as overrides over the
// default configuration.
func (cfg *Cfg) New() Cfg {
	return DefaultCfg().Override(cfg)
}

// Override merges other into cfg; non-zero fields from other override the
// corresponding values in cfg. If other is nil, the receiver is returned.
//
//goland:noinspection DuplicatedCode
func (cfg *Cfg) Override(other *Cfg) Cfg {
	if other == nil {
		return *cfg
	}
	result := Cfg{
		Name:              util.Either(other.Name, cfg.Name),
		Source:            util.Either(other.Source, cfg.Source),
		Start:             util.Either(other.Start, cfg.Start),
		NoMemo:            util.Either(other.NoMemo, cfg.NoMemo),
		NoPruneMemosOnCut: util.Either(other.NoPruneMemosOnCut, cfg.NoPruneMemosOnCut),
		PerLineMemos:      util.Either(other.PerLineMemos, cfg.PerLineMemos),
		Trace:             util.Either(other.Trace, cfg.Trace),
		Colorize:          util.Either(other.Colorize, cfg.Colorize),
		Grammar:           util.Either(other.Grammar, cfg.Grammar),
		NoLeftRecursion:   util.Either(other.NoLeftRecursion, cfg.NoLeftRecursion),
		IgnoreCase:        util.Either(other.IgnoreCase, cfg.IgnoreCase),
		NameChars:         util.Either(other.NameChars, cfg.NameChars),
		Whitespace:        util.Either(other.Whitespace, cfg.Whitespace),
		Comments:          util.Either(other.Comments, cfg.Comments),
		EolComments:       util.Either(other.EolComments, cfg.EolComments),
		Keywords:          util.EitherSlice(other.Keywords, cfg.Keywords),
		ParseInfo:         util.Either(other.ParseInfo, cfg.ParseInfo),
		Heart:             util.Either(other.Heart, cfg.Heart),
	}
	if other.Semantics != nil {
		result.Semantics = other.Semantics
	}

	if other.Grammar != "" {
		result.Name = result.Grammar
	}

	if result.IgnoreCase && len(result.Keywords) > 0 {
		upper := make([]string, len(result.Keywords))
		for i, kw := range result.Keywords {
			upper[i] = strings.ToUpper(kw)
		}
		result.Keywords = upper
	}

	if result.NoMemo {
		result.NoLeftRecursion = true
	}

	if result.NameChars != "" {
		result.NameGuard = new(true)
	}
	if other.NameGuard != nil {
		result.NameGuard = other.NameGuard
	}

	return result
}
