package context

import (
	"fmt"
	"path/filepath"
)

// Location represents a source code location with file and line number.
type Location struct {
	File string
	Line int
}

// Nope represents a localized parsing failure with a source location.
type Nope struct {
	location Location
}

// DisasterReport aggregates failure information (furthest failure), cut
// state, inner error and tracing memento for diagnostic reporting.
type DisasterReport struct {
	CutSeen  bool
	location Location
	Inner    error
	Memento  Memento
}

// Error returns a string representation of the Nope error.
func (e *Nope) Error() string {
	return fmt.Sprintf("Scaped Nope at %v", e.location)
}

// Error returns a string representation of the DisasterReport.
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

// Start returns the starting position of the failure.
func (e *DisasterReport) Start() int {
	return e.Memento.Start
}

// Mark returns the input mark at which the failure occurred.
func (e *DisasterReport) Mark() int {
	return e.Memento.Mark
}

// Unwrap returns the inner error of the DisasterReport.
func (e *DisasterReport) Unwrap() error {
	return e.Inner
}
