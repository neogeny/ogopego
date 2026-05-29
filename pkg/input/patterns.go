// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"github.com/neogeny/ogopego/pkg/config"
	pyre2 "github.com/neogeny/ogopego/pkg/util/pyre"
)

// TokenizingPatterns groups precompiled patterns used by cursors for
// whitespace, comments and EOL detection.
type TokenizingPatterns struct {
	Wsp        pyre2.Pattern // Wsp is the whitespace pattern.
	Cmt        pyre2.Pattern // Cmt is the comment pattern.
	Eol        pyre2.Pattern // Eol is the end-of-line comment pattern.
	NonDefault bool
}

func NewPatterns(wsp string, cmt string, eol string) (*TokenizingPatterns, error) {
	var err error
	var pwsp pyre2.Pattern
	var pcmt pyre2.Pattern
	var peol pyre2.Pattern

	if pwsp, err = pyre2.Compile(wsp); err != nil {
		return nil, err
	}
	if pcmt, err = pyre2.Compile(cmt); err != nil {
		return nil, err
	}
	if peol, err = pyre2.Compile(eol); err != nil {
		return nil, err
	}
	return &TokenizingPatterns{
		Wsp: pwsp,
		Cmt: pcmt,
		Eol: peol,
	}, nil

}

func DefaultPatterns() TokenizingPatterns {
	pat, err := NewPatterns(`(?m)\s+`, `(?m)#.*`, `(?m)#.*$`)
	if err != nil {
		panic("failed to compile default patterns: " + err.Error())
	}
	return *pat
}

func (p *TokenizingPatterns) Configure(cfg config.Cfg) {
	p.NonDefault = false
	if cfg.Whitespace != nil {
		p.NonDefault = true
		if *cfg.Whitespace != "" {
			if pat, err := pyre2.Compile(*cfg.Whitespace); err == nil {
				p.Wsp = pat
			}
		} else {
			p.Wsp = nil
		}
	}
	p.Cmt = nil
	if cfg.Comments != "" {
		p.NonDefault = true
		if pat, err := pyre2.Compile(cfg.Comments); err == nil {
			p.Cmt = pat
		}
	}

	p.Eol = nil
	if cfg.EolComments != "" {
		p.NonDefault = true
		if pat, err := pyre2.Compile(cfg.EolComments); err == nil {
			p.Eol = pat
		}
	}

}
