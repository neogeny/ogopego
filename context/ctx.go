package context

import (
	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/util/pyre"
)

type Ctx interface {
	Cursor() input.Cursor
	CursorMut() input.Cursor
	CallStack() CallStack
	CutSeen() bool
	Tracer() Tracer
	Mark() int
	Enter(name string)
	Leave()
	Track(key MemoKey) int
	Untrack(key MemoKey) int
	Failure(start int, source ParseFailure) *DisasterReport
	SetFurthestFailure(dis *DisasterReport)
	FurthestFailure() *DisasterReport
	Reset(mark int)
	AtEnd() bool
	ParseEOF() bool
	Dot() (rune, bool)
	Next() (rune, bool)
	Peek() (rune, bool)
	MatchEOL() bool
	MatchVoid()
	NextToken()
	GetPattern(pattern string) pyre.Pattern
	MatchToken(token string) bool
	MatchPattern(pattern string) (string, bool)
	Intern(s string) string
	Key(name string, canMemo bool, mark int) MemoKey
	Memo(key MemoKey) *Memo
	Memoize(key MemoKey, tree interface{}, lastMark int)
	ClearErrorMemos()
	PruneCache(cutpoint int)
	Cut()
	ClearCut()
	IsKeyword(name string) bool
	SetKeywords(keywords []string)
	Merge(other Ctx) Ctx
	Push() Ctx
	HeartbeatTick()
}

type BaseCtx struct {
	cursor       input.Cursor
	cursorMut    input.Cursor
	callStack    CallStack
	tracer       Tracer
	cut          bool
	furthest     *DisasterReport
	patternCache map[string]pyre.Pattern
	keywords     []string
}

func (b *BaseCtx) Cursor() input.Cursor {
	return b.cursor
}

func (b *BaseCtx) CursorMut() input.Cursor {
	return b.cursorMut
}

func (b *BaseCtx) CallStack() CallStack {
	return b.callStack
}

func (b *BaseCtx) CutSeen() bool {
	return b.cut
}

func (b *BaseCtx) Tracer() Tracer {
	return b.tracer
}

func (b *BaseCtx) Mark() int {
	return b.cursor.Mark()
}

func (b *BaseCtx) Reset(mark int) {
	b.cursorMut.Reset(mark)
}

func (b *BaseCtx) AtEnd() bool {
	return b.cursor.AtEnd()
}

func (b *BaseCtx) Next() (rune, bool) {
	return b.cursorMut.Next()
}

func (b *BaseCtx) Peek() (rune, bool) {
	return b.cursorMut.Peek()
}

func (b *BaseCtx) Dot() (rune, bool) {
	return b.Next()
}

func (b *BaseCtx) MatchEOL() bool {
	return b.cursorMut.MatchEOL()
}

func (b *BaseCtx) MatchVoid() {
	b.NextToken()
}

func (b *BaseCtx) NextToken() {
	b.cursorMut.NextToken()
}

func (b *BaseCtx) HeartbeatTick() {}

func (b *BaseCtx) Intern(s string) string {
	return s
}

func (b *BaseCtx) IsKeyword(name string) bool {
	for _, kw := range b.keywords {
		if kw == name {
			return true
		}
	}
	return false
}

func (b *BaseCtx) SetKeywords(keywords []string) {
	b.keywords = keywords
}

func (b *BaseCtx) Cut() {
	b.cut = true
}

func (b *BaseCtx) ClearCut() {
	b.cut = false
}

func (b *BaseCtx) SetFurthestFailure(dis *DisasterReport) {
	b.furthest = dis
}

func (b *BaseCtx) FurthestFailure() *DisasterReport {
	return b.furthest
}

func (b *BaseCtx) Failure(start int, source ParseFailure) *DisasterReport {
	b.cursorMut.Reset(start)

	if f := b.FurthestFailure(); f != nil && f.Start >= b.Mark() {
		f.CutSeen = b.cut
		return f
	}

	dis := &DisasterReport{Start: start, Failure: source, CutSeen: b.cut}
	b.SetFurthestFailure(dis)
	return dis
}

func (b *BaseCtx) GetPattern(pattern string) pyre.Pattern {
	if b.patternCache == nil {
		b.patternCache = make(map[string]pyre.Pattern)
	}
	if p, ok := b.patternCache[pattern]; ok {
		return p
	}
	p, err := pyre.Compile(pattern)
	if err != nil {
		return nil
	}
	b.patternCache[pattern] = p
	return p
}

func (b *BaseCtx) MatchToken(token string) bool {
	b.NextToken()
	result := b.cursorMut.PeekToken(token)
	if result {
		b.cursorMut.MatchToken(token)
	}
	return result
}

func (b *BaseCtx) MatchPattern(pattern string) (string, bool) {
	re := b.GetPattern(pattern)
	return b.cursorMut.MatchPattern(re)
}

func (b *BaseCtx) Key(name string, canMemo bool, mark int) MemoKey {
	return MemoKey{Mark: mark, Name: name, CanMemo: canMemo}
}

func (b *BaseCtx) Memo(key MemoKey) *Memo {
	return nil
}

func (b *BaseCtx) Memoize(key MemoKey, tree interface{}, lastMark int) {}

func (b *BaseCtx) ClearErrorMemos() {}

func (b *BaseCtx) PruneCache(cutpoint int) {}

func (b *BaseCtx) Enter(name string) {
	b.callStack = append(b.callStack, name)
}

func (b *BaseCtx) Leave() {
	if len(b.callStack) > 0 {
		b.callStack = b.callStack[:len(b.callStack)-1]
	}
}

func (b *BaseCtx) Track(key MemoKey) int {
	return 0
}

func (b *BaseCtx) Untrack(key MemoKey) int {
	return 0
}

func (b *BaseCtx) ParseEOF() bool {
	b.Enter("__eof__")
	b.NextToken()
	result := b.cursor.AtEnd()
	b.Leave()
	return result
}

func (b *BaseCtx) Push() Ctx {
	return nil
}

func (b *BaseCtx) Merge(other Ctx) Ctx {
	return nil
}
