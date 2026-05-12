package context

import (
	"fmt"
	"sort"

	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/tree"
	"github.com/neogeny/ogopego/util/pyre"
)

type BaseCtx struct {
	cursor       input.Cursor
	callStack    CallStack
	tracer       Tracer
	furthest     *DisasterReport
	patternCache map[string]pyre.Pattern
	keywords     map[string]struct{}
}

func NewBaseCtx(cursor input.Cursor) *BaseCtx {
	return &BaseCtx{
		cursor: cursor,
		tracer: NullTracer{},
	}
}

func NewBaseCtxWithTracer(cursor input.Cursor, tracer Tracer) *BaseCtx {
	return &BaseCtx{
		cursor: cursor,
		tracer: tracer,
	}
}

func (b *BaseCtx) Cursor() input.Cursor { return b.cursor }

func (b *BaseCtx) CallStack() CallStack { return b.callStack }

func (b *BaseCtx) Tracer() Tracer { return b.tracer }

func (b *BaseCtx) Mark() int { return b.cursor.Mark() }

func (b *BaseCtx) Reset(mark int) { b.cursor.Reset(mark) }

func (b *BaseCtx) AtEnd() bool { return b.cursor.AtEnd() }

func (b *BaseCtx) Next() (rune, bool) { return b.cursor.Next() }

func (b *BaseCtx) Peek() (rune, bool) { return b.cursor.Peek() }

func (b *BaseCtx) Dot() (rune, error) {
	r, ok := b.Next()
	if !ok {
		return 0, &ParseError{Pos: b.Mark(), Message: "expected any character"}
	}
	return r, nil
}

func (b *BaseCtx) MatchEOL() bool { return b.cursor.MatchEOL() }

func (b *BaseCtx) NextToken() { b.cursor.NextToken() }

func (b *BaseCtx) HeartbeatTick() {}

func (b *BaseCtx) Intern(s string) string { return s }

func (b *BaseCtx) IsKeyword(name string) bool {
	_, ok := b.keywords[name]
	return ok
}

func (b *BaseCtx) SetKeywords(keywords []string) {
	sort.Strings(keywords)
	set := make(map[string]struct{}, len(keywords))
	for _, kw := range keywords {
		set[kw] = struct{}{}
	}
	b.keywords = set
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
	return b.cursor.MatchToken(token)
}

func (b *BaseCtx) MatchPattern(pattern string) (string, bool) {
	re := b.GetPattern(pattern)
	return b.cursor.MatchPattern(re)
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
	msg := source.Error()
	dis := &DisasterReport{
		Start:   start,
		Failure: ParseFailure{Message: msg},
		Memento: input.NewMemento(start, msg, b.cursor, b.callStack),
	}
	b.SetFurthestFailure(dis)
	return dis
}

func (b *BaseCtx) FurthestFailure() *DisasterReport { return b.furthest }

func (b *BaseCtx) SetFurthestFailure(dis *DisasterReport) { b.furthest = dis }

func (b *BaseCtx) Token(token string) (string, error) {
	b.NextToken()
	ok := b.cursor.MatchToken(token)
	if !ok {
		return "", &ParseError{Pos: b.Mark(), Message: fmt.Sprintf("expected %q", token)}
	}
	return token, nil
}

func (b *BaseCtx) Pattern(pattern string) (string, error) {
	b.NextToken()
	re := b.GetPattern(pattern)
	if re == nil {
		return "", &ParseError{Pos: b.Mark(), Message: fmt.Sprintf("invalid pattern %q", pattern)}
	}
	m, ok := b.cursor.MatchPattern(re)
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
	if !b.cursor.MatchEOL() {
		return &ParseError{Pos: b.Mark(), Message: "expected end of line"}
	}
	return nil
}

func (b *BaseCtx) Eof() bool { return b.cursor.AtEnd() }

func (b *BaseCtx) Constant(literal any) (tree.Tree, error) {
	switch v := literal.(type) {
	case string:
		return &tree.Text{Value: v}, nil
	case float64:
		return &tree.Number{Value: v}, nil
	case bool:
		return &tree.Bool{Value: v}, nil
	case nil:
		return &tree.Nil{}, nil
	case int:
		return &tree.Number{Value: float64(v)}, nil
	default:
		return &tree.Text{Value: fmt.Sprintf("%v", v)}, nil
	}
}
