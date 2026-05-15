package context

import (
	"errors"
	"fmt"

	"github.com/neogeny/ogopego/input"
)

type Nope struct {
	error
	CutSeen bool
}

type ParseFailure struct {
	error
	Msg     string
	Start   int
	Mark    int
	CutSeen bool
	Inner   error
}

type DisasterReport struct {
	error
	CutSeen bool
	Inner   error
	Memento *input.Memento
}

func (e *Nope) Error() string {
	return fmt.Sprintf("Scaped Nope %v", e.CutSeen)
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
	if e.Memento != nil {
		return e.Memento.Error()
	}
	if e.Inner != nil {
		return e.Inner.Error()
	}
	return fmt.Sprintf("ParseError %v", e)
}

func (e *DisasterReport) Mark() int {
	return e.Memento.Mark
}

func (e *DisasterReport) Unwrap() error {
	return e.Inner
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
