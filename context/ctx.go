// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package context

import (
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

type Ctx interface {
	Configurable
	Cursor() Cursor
	CallStack() CallStack
	Tracer() Tracer
	Mark() int
	Reset(mark int)
	AtEnd() bool
	Next() (rune, bool)
	Peek() (rune, bool)
	Dot() (rune, error)
	NextToken()
	MatchEOL() bool
	MatchToken(token string) bool
	MatchPattern(pattern string) (string, bool)
	GetPattern(pattern string) pyre.Pattern
	Token(token string) (string, error)
	Pattern(pattern string) (string, error)
	Void() error
	Fail() error
	Eof() bool
	EofCheck() error
	EolCheck() error
	Constant(literal any) (trees.Tree, error)
	Enter(name string)
	Leave()
	Failure(start int, source error) *ParseFailure
	FurthestFailure() *DisasterReport
	SetFurthestFailure(dis *DisasterReport)
	IsKeyword(name string) bool
	Intern(s string) string
	ParseEOF() bool
	HeartbeatTick()
	Key(name string, canMemo bool) MemoKey
	Memo(key MemoKey) (Memo, bool)
	Memoize(key MemoKey, tree trees.Tree, mark int)
	TrackRecursionDepth(key MemoKey) error
	Untrack(key MemoKey)
}
