package context

import (
	"fmt"
	"sort"

	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

type BaseCtx struct {
	cfg            Cfg
	cursor         input.Cursor
	callStack      CallStack
	tracer         Tracer
	furthest       *DisasterReport
	patternCache   map[string]pyre.Pattern
	keywords       map[string]struct{}
	memoCache      map[MemoKey]Memo
	recursionKey   MemoKey
	recursionDepth int
}

func NewCtx(cursor input.Cursor, cfg *Cfg) *BaseCtx {
	ctx := BaseCtx{
		cfg:    cfg.New(),
		cursor: cursor,
		tracer: NullTracer{},
	}
	ctx.cursor.Configure(ctx.cfg)
	return &ctx
}

func (b *BaseCtx) Configure(cfg Cfg) {
	b.cursor.Configure(cfg)
	b.setKeywords(cfg.Keywords)
}

func (b *BaseCtx) SetTracer(tracer Tracer) {
	b.tracer = tracer
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
		return 0, &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("expected any character"),
		}
	}
	return r, nil
}

func (b *BaseCtx) MatchEOL() bool { return b.cursor.MatchEOL() }

func (b *BaseCtx) NextToken() { b.cursor.NextToken() }

func (b *BaseCtx) HeartbeatTick() {}

func (b *BaseCtx) Key(name string, canMemo bool) MemoKey {
	return MemoKey{Mark: b.Mark(), Name: name, CanMemo: canMemo}
}

func (b *BaseCtx) Memo(key MemoKey) (Memo, bool) {
	if b.memoCache == nil {
		return Memo{}, false
	}
	m, ok := b.memoCache[key]
	return m, ok
}

func (b *BaseCtx) Memoize(key MemoKey, tree trees.Tree, mark int) {
	if !key.CanMemo {
		return
	}
	if b.memoCache == nil {
		b.memoCache = make(map[MemoKey]Memo)
	}
	b.memoCache[key] = Memo{Tree: tree, Mark: mark}
}

func (b *BaseCtx) TrackRecursionDepth(key MemoKey) error {
	if key == b.recursionKey {
		b.recursionDepth++
	} else {
		b.recursionKey = key
		b.recursionDepth = 1
	}
	if b.recursionDepth > 64 {
		return fmt.Errorf("recursion depth exceeded")
	}
	return nil
}

func (b *BaseCtx) Untrack(key MemoKey) {
	if key == b.recursionKey {
		b.recursionDepth--
		if b.recursionDepth <= 0 {
			b.recursionKey = MemoKey{}
			b.recursionDepth = 0
		}
	}
}

func (b *BaseCtx) Intern(s string) string { return s }

func (b *BaseCtx) IsKeyword(name string) bool {
	_, ok := b.keywords[name]
	return ok
}

func (b *BaseCtx) setKeywords(keywords []string) {
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

func (b *BaseCtx) Failure(start int, source error) *ParseFailure {
	b.Reset(start)
	nope := &ParseFailure{
		Start: start,
		Mark:  b.Mark(),
		Inner: source,
	}
	if furthest := b.FurthestFailure(); furthest != nil && furthest.Start >= b.Mark() {
		return nope
	}
	msg := source.Error()
	dis := &DisasterReport{
		Start:   start,
		Failure: nope,
		Memento: input.NewMemento(start, msg, b.cursor, b.callStack),
	}
	b.SetFurthestFailure(dis)
	return nope
}

func (b *BaseCtx) FurthestFailure() *DisasterReport { return b.furthest }

func (b *BaseCtx) SetFurthestFailure(dis *DisasterReport) { b.furthest = dis }

func (b *BaseCtx) Token(token string) (string, error) {
	b.NextToken()
	ok := b.cursor.MatchToken(token)
	if !ok {
		return "", &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("expected %q", token),
		}
	}
	return token, nil
}

func (b *BaseCtx) Pattern(pattern string) (string, error) {
	b.NextToken()
	re := b.GetPattern(pattern)
	if re == nil {
		return "", &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("invalid pattern %q", pattern)}
	}
	m, ok := b.cursor.MatchPattern(re)
	if !ok {
		return "", &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("expected pattern %q", pattern),
		}
	}
	return m, nil
}

func (b *BaseCtx) Void() error {
	b.NextToken()
	return nil
}

func (b *BaseCtx) Fail() error {
	return &ParseFailure{
		Mark:  b.Mark(),
		Inner: fmt.Errorf("fail")}
}

func (b *BaseCtx) EofCheck() error {
	b.NextToken()
	if !b.cursor.AtEnd() {
		return &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("expected end of text"),
		}
	}
	return nil
}

func (b *BaseCtx) EolCheck() error {
	if !b.cursor.MatchEOL() {
		return &ParseFailure{
			Mark:  b.Mark(),
			Inner: fmt.Errorf("expected end of line"),
		}
	}
	return nil
}

func (b *BaseCtx) Eof() bool { return b.cursor.AtEnd() }

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
