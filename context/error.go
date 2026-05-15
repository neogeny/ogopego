package context

import (
	"errors"
	"fmt"
)

type Location struct {
	File string
	Line int
}

type Nope struct {
	error
	CutSeen  bool
	location Location
}

type DisasterReport struct {
	error
	CutSeen  bool
	location Location
	Inner    error
	Memento  Memento
}

func (e *Nope) Error() string {
	return fmt.Sprintf("Scaped Nope %v at [%v]", e.CutSeen, e.location)
}

func (e *DisasterReport) Error() string {
	var inner string
	if e.Inner != nil {
		inner = e.Inner.Error()
	}

	memento := e.Memento.Error()

	// This ensures go test prints the full structural context of the DisasterReport
	return fmt.Sprintf(
		"DisasterReport [CutSeen: %t, Loc: %v]: %s\n%s",
		e.CutSeen,
		e.location,
		inner,
		memento,
	)
}

func (e *DisasterReport) Start() int {
	return e.Memento.Start
}

func (e *DisasterReport) Mark() int {
	return e.Memento.Mark
}

func (e *DisasterReport) Unwrap() error {
	return e.Inner
}

func MarkCut(err error, value bool) error {
	var pf *Nope
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
	if pf, ok := errors.AsType[*Nope](err); ok && pf.CutSeen {
		pf.CutSeen = false
		return true
	}
	return false
}

func IsCut(err error) bool {
	if err == nil {
		return false
	}
	if pf, ok := errors.AsType[*Nope](err); ok && pf.CutSeen {
		return true
	}
	return false
}
