// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package config

import (
	"strings"

	"github.com/neogeny/ogopego/util/heartbeat"
)

const DefaultPerlinememos = 8

type Configurable interface {
	Configure(cfg Cfg)
}

// Cfg configures grammar compilation and input parsing. Use DefaultCfg() or
// pass nil to API functions for defaults. Individual fields override
// grammar-level @@directives.
type Cfg struct {
	Name   string // grammar name (overrides @@grammar)
	Source string // source description for error messages
	Start  string // start rule name (default: "start")

	Semantics any // reserved for future use

	NoMemo            bool    // disable memoization
	NoPruneMemosOnCut bool    // disable pruning memos on cut (~)
	PerLineMemos      float64 // memo entries per input line

	Trace    bool // enable parse trace output to stderr
	Colorize bool // colorize trace output

	Grammar         string // grammar text (overrides @@grammar body)
	NoLeftRecursion bool   // disable left-recursion support

	IgnoreCase bool   // case-insensitive matching
	NameChars  string // additional name characters (implies NameGuard)
	NameGuard  bool   // enforce word boundaries on names

	Whitespace  *string // whitespace pattern (nil = default, &"" = none)
	Comments    string  // comment pattern
	EolComments string  // end-of-line comment pattern

	Keywords []string // reserved words

	ParseInfo bool // attach source position info to AST nodes

	Heartbeat heartbeat.Heartbeat // progress callback (CLI progress bars)
}

func Either[T comparable](userVal, defaultVal T) T {
	var zero T
	if userVal != zero {
		return userVal
	}
	return defaultVal
}

func eitherSlice[T any](userVal, defaultVal []T) []T {
	if userVal != nil {
		return userVal
	}
	return defaultVal
}

// DefaultCfg returns the default configuration. Pass nil to API functions
// to use these defaults.
func DefaultCfg() Cfg {
	ws := `(?m)\s+`
	return Cfg{
		NoMemo:            false,
		NoPruneMemosOnCut: false,
		PerLineMemos:      DefaultPerlinememos,
		Trace:             false,
		Colorize:          false,
		NoLeftRecursion:   false,
		IgnoreCase:        false,
		NameGuard:         false,
		Whitespace:        &ws,
		Keywords:          nil,
		ParseInfo:         false,
	}
}

func (cfg *Cfg) New() Cfg {
	return DefaultCfg().Override(cfg)
}

// Override merges an optional cfg over the receiver. Fields set on other
// override the receiver's values. Returns the receiver if other is nil.
func (cfg Cfg) Override(other *Cfg) Cfg {
	if other == nil {
		return cfg
	}
	result := Cfg{
		Name:              Either(other.Name, cfg.Name),
		Source:            Either(other.Source, cfg.Source),
		Start:             Either(other.Start, cfg.Start),
		Semantics:         Either(other.Semantics, cfg.Semantics),
		NoMemo:            Either(other.NoMemo, cfg.NoMemo),
		NoPruneMemosOnCut: Either(other.NoPruneMemosOnCut, cfg.NoPruneMemosOnCut),
		PerLineMemos:      Either(other.PerLineMemos, cfg.PerLineMemos),
		Trace:             Either(other.Trace, cfg.Trace),
		Colorize:          Either(other.Colorize, cfg.Colorize),
		Grammar:           Either(other.Grammar, cfg.Grammar),
		NoLeftRecursion:   Either(other.NoLeftRecursion, cfg.NoLeftRecursion),
		IgnoreCase:        Either(other.IgnoreCase, cfg.IgnoreCase),
		NameChars:         Either(other.NameChars, cfg.NameChars),
		NameGuard:         Either(other.NameGuard, cfg.NameGuard),
		Whitespace:        Either(other.Whitespace, cfg.Whitespace),
		Comments:          Either(other.Comments, cfg.Comments),
		EolComments:       Either(other.EolComments, cfg.EolComments),
		Keywords:          eitherSlice(other.Keywords, cfg.Keywords),
		ParseInfo:         Either(other.ParseInfo, cfg.ParseInfo),
		Heartbeat:         Either(other.Heartbeat, cfg.Heartbeat),
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
		result.NameGuard = true
	}

	return result
}
