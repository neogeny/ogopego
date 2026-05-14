// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

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
	cursor         Cursor
	callStack      CallStack
	tracer         Tracer
	furthest       *DisasterReport
	patternCache   map[string]pyre.Pattern
	keywords       map[string]struct{}
	memoCache      map[MemoKey]Memo
	recursionKey   MemoKey
	recursionDepth int
}

func NewCtx(cursor Cursor, cfg *Cfg) *BaseCtx {
	ctx := BaseCtx{
		cfg:    cfg.New(),
		cursor: cursor,
		tracer: NullTracer{},
	}
	ctx.cursor.Configure(ctx.cfg)
	return &ctx
}

func (ctx *BaseCtx) Configure(cfg Cfg) {
	ctx.cursor.Configure(cfg)
	ctx.setKeywords(cfg.Keywords)
}

func (ctx *BaseCtx) SetTracer(tracer Tracer) {
	ctx.tracer = tracer
}

func (ctx *BaseCtx) Cursor() Cursor { return ctx.cursor }

func (ctx *BaseCtx) CallStack() CallStack { return ctx.callStack }

func (ctx *BaseCtx) Tracer() Tracer { return ctx.tracer }

func (ctx *BaseCtx) Mark() int { return ctx.cursor.Mark() }

func (ctx *BaseCtx) Reset(mark int) { ctx.cursor.Reset(mark) }

func (ctx *BaseCtx) AtEnd() bool { return ctx.cursor.AtEnd() }

func (ctx *BaseCtx) Next() (rune, bool) { return ctx.cursor.Next() }

func (ctx *BaseCtx) Peek() (rune, bool) { return ctx.cursor.Peek() }

func (ctx *BaseCtx) Dot() (rune, error) {
	r, ok := ctx.Next()
	if !ok {
		return 0, &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("expected any character"),
		}
	}
	return r, nil
}

func (ctx *BaseCtx) MatchEOL() bool { return ctx.cursor.MatchEOL() }

func (ctx *BaseCtx) NextToken() { ctx.cursor.NextToken() }

func (ctx *BaseCtx) HeartbeatTick() {}

func (ctx *BaseCtx) Key(name string, canMemo bool) MemoKey {
	return MemoKey{Mark: ctx.Mark(), Name: name, CanMemo: canMemo}
}

func (ctx *BaseCtx) Memo(key MemoKey) (Memo, bool) {
	if ctx.memoCache == nil {
		return Memo{}, false
	}
	m, ok := ctx.memoCache[key]
	return m, ok
}

func (ctx *BaseCtx) Memoize(key MemoKey, tree trees.Tree, mark int) {
	if !key.CanMemo {
		return
	}
	if ctx.memoCache == nil {
		ctx.memoCache = make(map[MemoKey]Memo)
	}
	ctx.memoCache[key] = Memo{Tree: tree, Mark: mark}
}

func (ctx *BaseCtx) TrackRecursionDepth(key MemoKey) error {
	if key == ctx.recursionKey {
		ctx.recursionDepth++
	} else {
		ctx.recursionKey = key
		ctx.recursionDepth = 1
	}
	if ctx.recursionDepth > 64 {
		return fmt.Errorf("recursion depth exceeded")
	}
	return nil
}

func (ctx *BaseCtx) Untrack(key MemoKey) {
	if key == ctx.recursionKey {
		ctx.recursionDepth--
		if ctx.recursionDepth <= 0 {
			ctx.recursionKey = MemoKey{}
			ctx.recursionDepth = 0
		}
	}
}

func (ctx *BaseCtx) Intern(s string) string { return s }

func (ctx *BaseCtx) IsKeyword(name string) bool {
	_, ok := ctx.keywords[name]
	return ok
}

func (ctx *BaseCtx) setKeywords(keywords []string) {
	sort.Strings(keywords)
	set := make(map[string]struct{}, len(keywords))
	for _, kw := range keywords {
		set[kw] = struct{}{}
	}
	ctx.keywords = set
}

func (ctx *BaseCtx) GetPattern(pattern string) pyre.Pattern {
	if ctx.patternCache == nil {
		ctx.patternCache = make(map[string]pyre.Pattern)
	}
	if p, ok := ctx.patternCache[pattern]; ok {
		return p
	}
	p, err := pyre.Compile(pattern)
	if err != nil {
		return nil
	}
	ctx.patternCache[pattern] = p
	return p
}

func (ctx *BaseCtx) MatchToken(token string) bool {
	ctx.NextToken()

	wordlike := true
	for _, r := range token {
		if !isAlphaNum(r) {
			wordlike = false
			break
		}
	}

	var result bool
	if wordlike && ctx.cursor.NameGuard() {
		var pat string
		if ctx.cursor.IgnoreCase() {
			pat = token + `\b`
		} else {
			pat = `(?i)` + token + `\b`
		}
		_, result = ctx.MatchPattern(pat)
	} else {
		result = ctx.cursor.MatchToken(token)
	}

	if result {
		ctx.Tracer().TraceMatch(ctx, token, "")
	} else {
		ctx.Tracer().TraceNoMatch(ctx, token, "")
	}
	return result
}

func isAlphaNum(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func (ctx *BaseCtx) MatchPattern(pattern string) (string, bool) {
	re := ctx.GetPattern(pattern)
	return ctx.cursor.MatchPattern(re)
}

func (ctx *BaseCtx) Enter(name string) {
	ctx.callStack = append(ctx.callStack, name)
}

func (ctx *BaseCtx) Leave() {
	if len(ctx.callStack) > 0 {
		ctx.callStack = ctx.callStack[:len(ctx.callStack)-1]
	}
}

func (ctx *BaseCtx) ParseEOF() bool {
	ctx.Enter("__eof__")
	ctx.NextToken()
	result := ctx.cursor.AtEnd()
	ctx.Leave()
	return result
}

func (ctx *BaseCtx) Failure(start int, source error) *ParseFailure {
	ctx.Reset(start)
	nope := &ParseFailure{
		Start: start,
		Mark:  ctx.Mark(),
		Inner: source,
	}
	if furthest := ctx.FurthestFailure(); furthest != nil && furthest.Start >= ctx.Mark() {
		return nope
	}
	msg := source.Error()
	dis := &DisasterReport{
		Start:   start,
		Failure: nope,
		Memento: input.NewMemento(start, msg, ctx.cursor, ctx.callStack),
	}
	ctx.SetFurthestFailure(dis)
	return nope
}

func (ctx *BaseCtx) FurthestFailure() *DisasterReport { return ctx.furthest }

func (ctx *BaseCtx) SetFurthestFailure(dis *DisasterReport) { ctx.furthest = dis }

func (ctx *BaseCtx) Token(token string) (string, error) {
	ctx.NextToken()
	ok := ctx.cursor.MatchToken(token)
	if !ok {
		return "", &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("expected %q", token),
		}
	}
	return token, nil
}

func (ctx *BaseCtx) Pattern(pattern string) (string, error) {
	ctx.NextToken()
	re := ctx.GetPattern(pattern)
	if re == nil {
		return "", &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("invalid pattern %q", pattern)}
	}
	m, ok := ctx.cursor.MatchPattern(re)
	if !ok {
		return "", &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("expected pattern %q", pattern),
		}
	}
	return m, nil
}

func (ctx *BaseCtx) Void() error {
	ctx.NextToken()
	return nil
}

func (ctx *BaseCtx) Fail() error {
	return &ParseFailure{
		Mark:  ctx.Mark(),
		Inner: fmt.Errorf("fail")}
}

func (ctx *BaseCtx) EofCheck() error {
	ctx.NextToken()
	if !ctx.cursor.AtEnd() {
		return &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("expected end of text"),
		}
	}
	return nil
}

func (ctx *BaseCtx) EolCheck() error {
	if !ctx.cursor.MatchEOL() {
		return &ParseFailure{
			Mark:  ctx.Mark(),
			Inner: fmt.Errorf("expected end of line"),
		}
	}
	return nil
}

func (ctx *BaseCtx) Eof() bool { return ctx.cursor.AtEnd() }

func (ctx *BaseCtx) Constant(literal any) (trees.Tree, error) {
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
