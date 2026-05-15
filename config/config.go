package config

import "strings"

const DefaultPerlinememos = 8

type Configurable interface {
	Configure(cfg Cfg)
}

type Cfg struct {
	Name   string
	Source string
	Start  string

	Semantics any

	NoMemo            bool
	NoPruneMemosOnCut bool
	PerLineMemos      float64

	Trace bool

	Grammar         string
	NoLeftRecursion bool

	IgnoreCase bool
	NameChars  string
	NameGuard  bool

	Whitespace  *string
	Comments    string
	EolComments string

	Keywords []string

	ParseInfo bool
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

func DefaultCfg() Cfg {
	ws := `(?m)\s+`
	return Cfg{
		NoMemo:            false,
		NoPruneMemosOnCut: false,
		PerLineMemos:      DefaultPerlinememos,
		Trace:             false,
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
