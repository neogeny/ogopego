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

// ParseFailure aggregates failure information (furthest failure), cut
// state, inner error and tracing memento for diagnostic reporting.
type ParseFailure struct {
	Inner    error
	location Location
	Memento  Memento
}

// Error returns a string representation of the ParseFailure.
func (e *ParseFailure) Error() string {
	var inner string
	if e.Inner != nil {
		inner = e.Inner.Error()
	}

	memento := e.Memento.Error()

	filename := filepath.Base(e.location.File)
	location := fmt.Sprintf("%s@%d", filename, e.location.Line)

	return fmt.Sprintf(
		"\nDisasterReport(\nLoc: %v,\nError: %s\n\n%s\n)",
		location,
		inner,
		memento,
	)
}

// Start returns the starting position of the failure.
func (e *ParseFailure) Start() int {
	return e.Memento.Start
}

// Mark returns the input mark at which the failure occurred.
func (e *ParseFailure) Mark() int {
	return e.Memento.Mark
}

// Unwrap returns the inner error of the ParseFailure.
func (e *ParseFailure) Unwrap() error {
	return e.Inner
}
