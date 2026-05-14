// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package context

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/input"
)

type ParseFailure struct {
	Start   int
	Mark    int
	CutSeen bool
	Inner   error
}

type DisasterReport struct {
	Start   int
	CutSeen bool
	Failure *ParseFailure
	Memento *input.Memento
}

func (e *ParseFailure) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("at %d: %v", e.Mark, e.Inner)
	}
	return fmt.Sprintf("at %d: ParseError", e.Mark)
}

func (e *ParseFailure) Unwrap() error {
	return e.Inner
}

func (e *DisasterReport) Error() string {
	if e.Failure != nil {
		return e.Failure.Error()
	}
	return fmt.Sprintf("at %d: ParseError", e.Failure.Mark)
}

func (e *DisasterReport) Unwrap() error {
	return e.Failure
}

type NoMatch struct {
	Pos     int
	Message string
}

func (e *NoMatch) Error() string {
	return fmt.Sprintf("at %d: %s", e.Pos, e.Message)
}

var ErrNoMatch = &NoMatch{Pos: -1, Message: "no match"}

func MarkCut(err error, value bool) error {
	var pf *ParseFailure
	ok := errors.As(err, &pf)
	if !ok {
		return err
	}
	pf.CutSeen = pf.CutSeen || value
	return pf
}

func TakeCut(err error) bool {
	if err == nil {
		return false
	}
	if pf, ok := errors.AsType[*ParseFailure](err); ok && pf.CutSeen {
		pf.CutSeen = false
		return true
	}
	return false
}

func IsCut(err error) bool {
	if err == nil {
		return false
	}
	if pf, ok := errors.AsType[*ParseFailure](err); ok && pf.CutSeen {
		return true
	}
	return false
}

func IsNoMatch(err error) bool {
	if err == nil {
		return false
	}
	var nm *NoMatch
	return errors.As(err, &nm)
}
