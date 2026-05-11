package context

import (
	"errors"
	"fmt"
)

var ErrCut = errors.New("cut")

type ParseError struct {
	Pos     int
	Message string
	Inner   error
}

func (e *ParseError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("at %d: %s: %v", e.Pos, e.Message, e.Inner)
	}
	return fmt.Sprintf("at %d: %s", e.Pos, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Inner
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
