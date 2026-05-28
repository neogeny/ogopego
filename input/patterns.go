// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/util/pyre"
)

// TokenizingPatterns groups precompiled patterns used by cursors for
// whitespace, comments and EOL detection.
type TokenizingPatterns struct {
	Wsp        pyre.Pattern // Wsp is the whitespace pattern.
	Cmt        pyre.Pattern // Cmt is the comment pattern.
	Eol        pyre.Pattern // Eol is the end-of-line comment pattern.
	NonDefault bool
}

func NewPatterns(wsp string, cmt string, eol string) (*TokenizingPatterns, error) {
	var err error
	var pwsp pyre.Pattern
	var pcmt pyre.Pattern
	var peol pyre.Pattern

	if pwsp, err = pyre.Compile(wsp); err != nil {
		return nil, err
	}
	if pcmt, err = pyre.Compile(cmt); err != nil {
		return nil, err
	}
	if peol, err = pyre.Compile(eol); err != nil {
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
			if pat, err := pyre.Compile(*cfg.Whitespace); err == nil {
				p.Wsp = pat
			}
		} else {
			p.Wsp = nil
		}
	}
	p.Cmt = nil
	if cfg.Comments != "" {
		p.NonDefault = true
		if pat, err := pyre.Compile(cfg.Comments); err == nil {
			p.Cmt = pat
		}
	}

	p.Eol = nil
	if cfg.EolComments != "" {
		p.NonDefault = true
		if pat, err := pyre.Compile(cfg.EolComments); err == nil {
			p.Eol = pat
		}
	}

}
