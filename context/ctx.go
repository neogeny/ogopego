package context

import (
	"fmt"

	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

type Ctx interface {
	Cursor() input.Cursor
	CursorMut() input.Cursor
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
	Failure(start int, source error) *DisasterReport
	FurthestFailure() *DisasterReport
	SetFurthestFailure(dis *DisasterReport)
	IsKeyword(name string) bool
	SetKeywords(keywords []string)
	Intern(s string) string
	ParseEOF() bool
	HeartbeatTick()
}

type BaseCtx struct {
	cursor       input.Cursor
	cursorMut    input.Cursor
	callStack    CallStack
	tracer       Tracer
	furthest     *DisasterReport
	patternCache map[string]pyre.Pattern
	keywords     []string
}

func NewBaseCtx(cursor input.Cursor) *BaseCtx {
	return &BaseCtx{
		cursor:    cursor,
		cursorMut: cursor,
	}
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

func (b *BaseCtx) Dot() (rune, error) {
	r, ok := b.Next()
	if !ok {
		return 0, &ParseError{Pos: b.Mark(), Message: "expected any character"}
	}
	return r, nil
}

func (b *BaseCtx) MatchEOL() bool {
	return b.cursorMut.MatchEOL()
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

func (b *BaseCtx) Enter(name string) {
	b.callStack = append(b.callStack, name)
}

func (b *BaseCtx) Leave() {
	if len(b.callStack) > 0 {
		b.callStack = b.callStack[:len(b.callStack)-1]
	}
}

func (b *BaseCtx) ParseEOF() bool {
	b.Enter("__eof__")
	b.NextToken()
	result := b.cursor.AtEnd()
	b.Leave()
	return result
}

func (b *BaseCtx) Failure(start int, source error) *DisasterReport {
	b.Reset(start)
	if f := b.FurthestFailure(); f != nil && f.Start >= b.Mark() {
		return f
	}
	dis := &DisasterReport{Start: start, Failure: ParseFailure{Message: source.Error()}}
	b.SetFurthestFailure(dis)
	return dis
}

func (b *BaseCtx) FurthestFailure() *DisasterReport {
	return b.furthest
}

func (b *BaseCtx) SetFurthestFailure(dis *DisasterReport) {
	b.furthest = dis
}

func (b *BaseCtx) Token(token string) (string, error) {
	b.NextToken()
	ok := b.cursorMut.MatchToken(token)
	if !ok {
		return "", &ParseError{Pos: b.Mark(), Message: fmt.Sprintf("expected %q", token)}
	}
	return token, nil
}

func (b *BaseCtx) Pattern(pattern string) (string, error) {
	re := b.GetPattern(pattern)
	if re == nil {
		return "", &ParseError{Pos: b.Mark(), Message: fmt.Sprintf("invalid pattern %q", pattern)}
	}
	m, ok := b.cursorMut.MatchPattern(re)
	if !ok {
		return "", &ParseError{Pos: b.Mark(), Message: fmt.Sprintf("expected pattern %q", pattern)}
	}
	return m, nil
}

func (b *BaseCtx) Void() error {
	b.NextToken()
	return nil
}

func (b *BaseCtx) Fail() error {
	return &ParseError{Pos: b.Mark(), Message: "fail"}
}

func (b *BaseCtx) EofCheck() error {
	b.NextToken()
	if !b.cursor.AtEnd() {
		return &ParseError{Pos: b.Mark(), Message: "expected end of text"}
	}
	return nil
}

func (b *BaseCtx) EolCheck() error {
	if !b.cursorMut.MatchEOL() {
		return &ParseError{Pos: b.Mark(), Message: "expected end of line"}
	}
	return nil
}

func (b *BaseCtx) Eof() bool {
	return b.cursor.AtEnd()
}

func (b *BaseCtx) Constant(literal any) (trees.Tree, error) {
	switch v := literal.(type) {
	case string:
		return &trees.Text{Value: v}, nil
	case float64:
		return &trees.Number{Value: v}, nil
	case bool:
		return &trees.Bool{Value: v}, nil
	case nil:
		return &trees.Nil{}, nil
	case int:
		return &trees.Number{Value: float64(v)}, nil
	default:
		return &trees.Text{Value: fmt.Sprintf("%v", v)}, nil
	}
}
