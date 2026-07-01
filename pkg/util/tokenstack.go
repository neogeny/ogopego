// Package token provides a simple Lisp-like List implementation.
package util

import (
	"iter"
	"strings"
)

// node defines the private interface for our list variants.
type node interface {
	newWithTail(string) node
}

type nilNode struct{}

func (n *nilNode) newWithTail(a string) node { return &atomNode{value: a} }

type atomNode struct {
	value string
}

func (n *atomNode) newWithTail(a string) node {
	return &consNode{car: n, cdr: &atomNode{value: a}}
}

type consNode struct {
	car node
	cdr node
}

func (n *consNode) newWithTail(a string) node {
	return &consNode{car: n.car, cdr: n.cdr.newWithTail(a)}
}

// TokenStack handles a Lisp-like cons-list of tokens, optimized for push (O(1)) and iteration.
type TokenStack struct {
	root node
}

// New creates an empty token stack.
func NewTokenStack() TokenStack {
	return TokenStack{root: &nilNode{}}
}

// Atom creates a token stack containing a single atom.
func Atom(a string) TokenStack {
	return TokenStack{root: &atomNode{value: a}}
}

// IsEmpty returns true if the stack contains no elements.
func (ts TokenStack) IsEmpty() bool {
	_, ok := ts.root.(*nilNode)
	return ok
}

// Push prepends a new atom to the list (O(1)).
func (ts *TokenStack) Push(a string) {
	atom := &atomNode{value: a}
	ts.root = &consNode{car: atom, cdr: ts.root}
}

// NewWithTail returns a new TokenStack with the provided string as the last element (O(N)).
func (ts TokenStack) NewWithTail(a string) TokenStack {
	return TokenStack{root: ts.root.newWithTail(a)}
}

// Tail returns the tail of the list (everything after the first element), if any.
func (ts TokenStack) Tail() (TokenStack, bool) {
	if cons, ok := ts.root.(*consNode); ok {
		return TokenStack{root: cons.cdr}, true
	}
	return TokenStack{}, false
}

// First returns the first element of the list, if any.
func (ts TokenStack) First() (string, bool) {
	current := ts.root
	for {
		switch n := current.(type) {
		case *atomNode:
			return n.value, true
		case *consNode:
			current = n.car
		case *nilNode:
			return "", false
		}
	}
}

// Len returns the total number of atoms currently contained within the stack (O(N)).
func (ts TokenStack) Len() int {
	var count int
	current := ts.root

	for {
		switch n := current.(type) {
		case *atomNode:
			return count + 1
		case *consNode:
			// In our layout, car is always an atomNode or a nested branch.
			// We walk down the car branch to count its atoms, then advance down cdr.
			count += TokenStack{root: n.car}.Len()
			current = n.cdr
		default: // *nilNode
			return count
		}
	}
}

// All returns a standard functional iterator over the elements in the stack.
// It matches the original Rust semantics exactly via recursive traversal.
func (ts TokenStack) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		var walk func(n node) bool
		walk = func(n node) bool {
			switch nodeType := n.(type) {
			case *atomNode:
				return yield(nodeType.value)
			case *consNode:
				// Process the left branch (car), then the right branch (cdr)
				if !walk(nodeType.car) {
					return false
				}
				return walk(nodeType.cdr)
			default: // *nilNode
				return true
			}
		}
		walk(ts.root)
	}
}

// ToSlice collects the stack elements into a slice of strings using the iterator.
func (ts TokenStack) ToSlice() []string {
	var items []string
	for val := range ts.All() {
		items = append(items, val)
	}
	return items
}

// String implements the fmt.Stringer interface.
func (ts TokenStack) String() string {
	return "[" + strings.Join(ts.ToSlice(), " ") + "]"
}
