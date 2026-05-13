package context

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/input"
)

var ErrCut = errors.New("cut")

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

type CutError struct {
	Err   error
	Start int
}

func (e *CutError) Error() string {
	return fmt.Sprintf("cut after %d: %v", e.Start, e.Err)
}

func (e *CutError) Unwrap() error {
	return e.Err
}

type NoMatch struct {
	Pos     int
	Message string
}

func (e *NoMatch) Error() string {
	return fmt.Sprintf("at %d: %s", e.Pos, e.Message)
}

var ErrNoMatch = &NoMatch{Pos: -1, Message: "no match"}

func IsCut(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrCut) || isCutType(err)
}

func isCutType(err error) bool {
	var ce *CutError
	return errors.As(err, &ce)
}

func IsNoMatch(err error) bool {
	if err == nil {
		return false
	}
	var nm *NoMatch
	return errors.As(err, &nm)
}
