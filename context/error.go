package context

import (
	"errors"
	"fmt"
	"path/filepath"
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

	filename := filepath.Base(e.location.File)
	location := fmt.Sprintf("%s@%d", filename, e.location.Line)

	return fmt.Sprintf(
		"\nDisasterReport(\n  CutSeen: %t,\n  Loc: %v,\n  Error: %s,\n%s\n)",
		e.CutSeen,
		location,
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
