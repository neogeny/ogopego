// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0
package context

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/pkg/util"
	"github.com/neogeny/ogopego/pkg/util/heartbeat"
)

// CallStack is a slice of call-site names representing the parser call stack.
type CallStack = util.TokenStack

// CoreCtxHeavy holds shared heavyweight state used across context clones.
type CoreCtxHeavy struct {
	mu            sync.Mutex
	cfg           Cfg
	memoCache     MemoCache
	tracer        Tracer
	keywords      map[string]struct{}
	heartbeat     heartbeat.Heart
	heartbeatTime time.Time
}

// CoreCtx is the concrete implementation of Ctx used by the parser runtime.
type CoreCtx struct {
	cursor         Cursor
	callStack      CallStack
	cutStack       []bool
	recursionKey   MemoKey
	recursionDepth int
	lookaheadDepth int
	lastCutMark    int
	furthest       *ParseFailure
	heavy          *CoreCtxHeavy
}

// NewCtx creates a new CoreCtx backed by the provided Cursor and optional
// configuration. Use the returned value where a context implementing Ctx is
// required.
func NewCtx(cursor Cursor, cfg *Cfg) *CoreCtx {
	cfgS := cfg.New()
	cursor.Configure(cfgS)

	stackCapacity := 64
	memoCapacity := max(
		stackCapacity,
		int(math.Round(cfgS.PerLineMemos*float64(cursor.LineCount()))),
	)
	heavy := &CoreCtxHeavy{
		cfg:       cfgS,
		memoCache: NewMemoMache(memoCapacity),
		tracer:    NullTracer{},
		heartbeat: heartbeat.NullHeart{},
	}
	ctx := CoreCtx{
		heavy:     heavy,
		cursor:    cursor,
		callStack: util.NewTokenStack(),
		cutStack:  make([]bool, 1, stackCapacity),
	}
	return &ctx
}

// Clone creates a deep copy of the CoreCtx, sharing the heavy state.
func (ctx *CoreCtx) Clone() Ctx {
	return &CoreCtx{
		cursor:         ctx.cursor.Clone(),
		callStack:      ctx.callStack,
		cutStack:       append([]bool(nil), ctx.cutStack...),
		recursionKey:   ctx.recursionKey,
		recursionDepth: ctx.recursionDepth,
		lookaheadDepth: ctx.lookaheadDepth,
		lastCutMark:    ctx.lastCutMark,
		heavy:          ctx.heavy,
	}
}

func (ctx *CoreCtx) Merge(other Ctx) {
	ctx.cursor = other.Cursor()
	ctx.furthest = other.FurthestFailure()
}

func (ctx *CoreCtx) muLock() {
	if ctx.heavy.cfg.Concurrency {
		ctx.heavy.mu.Lock()
	}
}

func (ctx *CoreCtx) muUnlock() {
	if ctx.heavy.cfg.Concurrency {
		ctx.heavy.mu.Unlock()
	}
}

func (ctx *CoreCtx) Cfg() Cfg {
	ctx.muLock()
	defer ctx.muUnlock()
	return ctx.heavy.cfg
}

func (ctx *CoreCtx) Configure(cfg Cfg) {
	ctx.muLock()
	ctx.heavy.cfg = ctx.heavy.cfg.Override(&cfg)
	ctx.muUnlock()
	ctx.cursor.Configure(cfg)

	ctx.muLock()
	ctx.setKeywords(cfg.Keywords)

	if cfg.Trace {
		if cfg.Colorize {
			color.Output = os.Stderr
			color.NoColor = false
		}
		ctx.heavy.tracer = ConsoleTracer{}
	} else {
		ctx.heavy.tracer = NullTracer{}
	}
	if cfg.Heart != nil {
		ctx.heavy.heartbeat = cfg.Heart
	}
	ctx.muUnlock()
}

func (ctx *CoreCtx) SetTracer(tracer Tracer) {
	ctx.muLock()
	ctx.heavy.tracer = tracer
	ctx.muUnlock()
}

func (ctx *CoreCtx) Cursor() Cursor { return ctx.cursor }

func (ctx *CoreCtx) CallStack() CallStack {
	return ctx.callStack
}

func (ctx *CoreCtx) Tracer() Tracer {
	ctx.muLock()
	defer ctx.muUnlock()
	return ctx.heavy.tracer
}

func (ctx *CoreCtx) Mark() int { return ctx.cursor.Mark() }

func (ctx *CoreCtx) Reset(mark int) { ctx.cursor.Reset(mark) }

func (ctx *CoreCtx) AtEnd() bool { return ctx.cursor.AtEnd() }

func (ctx *CoreCtx) Next() (rune, bool) { return ctx.cursor.Next() }

func (ctx *CoreCtx) Peek() (rune, bool) { return ctx.cursor.Peek() }

func (ctx *CoreCtx) Dot() (rune, error) {
	mark := ctx.Mark()
	r, ok := ctx.Next()
	if !ok {
		return 0, ctx.Failure(
			mark,
			fmt.Errorf("expected any character"),
		)
	}
	return r, nil
}

func (ctx *CoreCtx) MatchEOL() bool { return ctx.cursor.MatchEOL() }

func (ctx *CoreCtx) NextToken() { ctx.cursor.NextToken() }

func (ctx *CoreCtx) HeartbeatTick() {
	ctx.muLock()
	defer ctx.muUnlock()

	if time.Since(ctx.heavy.heartbeatTime) < 128*time.Millisecond {
		return
	}
	mark := ctx.Mark()
	total := ctx.cursor.Len()
	if total == 0 {
		return
	}
	ctx.heavy.heartbeat.Beat(mark, total)
	ctx.heavy.heartbeatTime = time.Now()
}

func (ctx *CoreCtx) pruneCache(cutPoint int) {
	ctx.muLock()
	defer ctx.muUnlock()
	PruneMemoCache(ctx.heavy.memoCache, cutPoint)
}

func (ctx *CoreCtx) Key(name string, canMemo bool) MemoKey {
	return MemoKey{Mark: ctx.Mark(), Name: name, CanMemo: canMemo}
}

func (ctx *CoreCtx) Memo(key MemoKey) (Memo, bool) {
	ctx.muLock()
	defer ctx.muUnlock()
	m, err := ctx.heavy.memoCache.Get(key)
	if err != nil {
		var zero Memo
		return zero, false
	}
	return m, true
}

func (ctx *CoreCtx) Memoize(key MemoKey, tree any, mark int) {
	if !key.CanMemo {
		return
	}
	ctx.muLock()
	_ = ctx.heavy.memoCache.Set(key, Memo{Tree: tree, Mark: mark})
	ctx.muUnlock()
}

func (ctx *CoreCtx) TrackRecursionDepth(key MemoKey) error {
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

func (ctx *CoreCtx) Untrack(key MemoKey) {
	if key == ctx.recursionKey {
		ctx.recursionDepth--
		if ctx.recursionDepth <= 0 {
			ctx.recursionKey = MemoKey{}
			ctx.recursionDepth = 0
		}
	}
}

func (ctx *CoreCtx) Intern(s string) string { return s }

func (ctx *CoreCtx) IsKeyword(name string) bool {
	ctx.muLock()
	defer ctx.muUnlock()
	_, ok := ctx.heavy.keywords[name]
	return ok
}

func (ctx *CoreCtx) setKeywords(keywords []string) {
	sort.Strings(keywords)
	set := make(map[string]struct{}, len(keywords))
	for _, kw := range keywords {
		set[kw] = struct{}{}
	}
	ctx.heavy.keywords = set
}

func (ctx *CoreCtx) MatchToken(token string) bool {
	ctx.NextToken()

	result := ctx.cursor.MatchToken(token)

	if result {
		ctx.Tracer().TraceMatch(ctx, token, "")
	} else {
		ctx.Tracer().TraceNoMatch(ctx, token, "")
	}
	return result
}

func (ctx *CoreCtx) Enter(name string) {
	ctx.callStack.Push(name)
}

func (ctx *CoreCtx) Leave() {
	if tail, ok := ctx.callStack.Tail(); ok {
		ctx.callStack = tail
	}
}

func (ctx *CoreCtx) ParseEOF() bool {
	ctx.Enter("__eof__")
	ctx.NextToken()
	result := ctx.cursor.AtEnd()
	ctx.Leave()
	return result
}

func (ctx *CoreCtx) Failure(start int, source error) error {
	loc := Location{}
	_, loc.File, loc.Line, _ = runtime.Caller(2)

	ctx.Reset(start)
	msg := source.Error()
	dis := &ParseFailure{
		location: loc,
		Inner:    source,
		Memento: NewMemento(
			start,
			msg,
			ctx.cursor,
			ctx.CallStack(),
		),
	}
	if ctx.furthest == nil || ctx.furthest.Mark() <= ctx.Mark() {
		ctx.SetFurthestFailure(dis)
	}
	return dis
}

func (ctx *CoreCtx) FurthestFailure() *ParseFailure {
	return ctx.furthest
}

func (ctx *CoreCtx) SetFurthestFailure(dis *ParseFailure) {
	ctx.furthest = dis
}

func (ctx *CoreCtx) MatchPattern(pattern string) (string, error) {
	mark := ctx.Mark()
	m, ok := ctx.cursor.MatchPattern(pattern)
	if !ok {
		return "", ctx.Failure(mark, fmt.Errorf("expected pattern %q", pattern))
	}
	return m, nil
}

func (ctx *CoreCtx) MatchName() (string, error) {
	ctx.NextToken()
	mark := ctx.Mark()
	m, ok := ctx.cursor.MatchName()
	if !ok {
		return "", ctx.Failure(mark, fmt.Errorf("expected name"))
	}
	return m, nil
}

func (ctx *CoreCtx) MatchInt() (int, error) {
	ctx.NextToken()
	mark := ctx.Mark()
	n, ok := ctx.cursor.MatchInt()
	if !ok {
		return 0, ctx.Failure(mark, fmt.Errorf("expected integer"))
	}
	return n, nil
}

func (ctx *CoreCtx) MatchUInt() (uint64, error) {
	ctx.NextToken()
	mark := ctx.Mark()
	n, ok := ctx.cursor.MatchUInt()
	if !ok {
		return 0, ctx.Failure(mark, fmt.Errorf("expected unsigned integer"))
	}
	return n, nil
}

func (ctx *CoreCtx) MatchFloat() (float64, error) {
	ctx.NextToken()
	mark := ctx.Mark()
	f, ok := ctx.cursor.MatchFloat()
	if !ok {
		return 0, ctx.Failure(mark, fmt.Errorf("expected float"))
	}
	return f, nil
}

func (ctx *CoreCtx) MatchBool() (bool, error) {
	ctx.NextToken()
	mark := ctx.Mark()
	b, ok := ctx.cursor.MatchBool()
	if !ok {
		return false, ctx.Failure(mark, fmt.Errorf("expected boolean"))
	}
	return b, nil
}

func (ctx *CoreCtx) Void() {
	ctx.NextToken()
}

func (ctx *CoreCtx) InLookahead() bool {
	return ctx.lookaheadDepth > 0
}

func (ctx *CoreCtx) EnterLookahead() {
	ctx.lookaheadDepth++
}

func (ctx *CoreCtx) LeaveLookahead() {
	ctx.lookaheadDepth--
}

func (ctx *CoreCtx) Fail() error {
	return ctx.Failure(
		ctx.Mark(),
		fmt.Errorf("fail"),
	)
}

func (ctx *CoreCtx) EofCheck() error {
	mark := ctx.Mark()
	ctx.NextToken()
	if !ctx.cursor.AtEnd() {
		return ctx.Failure(
			mark,
			fmt.Errorf("expected end of text"),
		)
	}
	return nil
}

func (ctx *CoreCtx) EolCheck() error {
	mark := ctx.Mark()
	if !ctx.cursor.MatchEOL() {
		return ctx.Failure(
			mark,
			fmt.Errorf("expected end of line"),
		)
	}
	return nil
}

func (ctx *CoreCtx) Eof() bool { return ctx.cursor.AtEnd() }

func (ctx *CoreCtx) Constant(literal any) (any, error) {
	switch v := literal.(type) {
	case string:
		return v, nil
	case float64:
		return v, nil
	case bool:
		return v, nil
	case nil:
		return nil, nil
	case int:
		return float64(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (ctx *CoreCtx) Cut() {
	ctx.cutStack[len(ctx.cutStack)-1] = true
	ctx.Tracer().TraceCut(ctx)

	if !ctx.heavy.cfg.NoPruneMemosOnCut && !ctx.InLookahead() {
		mark := ctx.Mark()
		if mark > ctx.lastCutMark {
			ctx.pruneCache(ctx.lastCutMark)
			ctx.lastCutMark = mark
		}
	}
}

func (ctx *CoreCtx) IsCutSeen() bool {
	return ctx.cutStack[len(ctx.cutStack)-1]
}
func (ctx *CoreCtx) CutStackPush() {
	ctx.cutStack = append(ctx.cutStack, false)
}

func (ctx *CoreCtx) CutStackPop() bool {
	cutSeen := ctx.IsCutSeen()
	ctx.cutStack = ctx.cutStack[:len(ctx.cutStack)-1]
	return cutSeen
}

func (ctx *CoreCtx) ApplySemantics(tree any, ruleName string, params []string) (any, bool) {
	ctx.muLock()
	sem := ctx.heavy.cfg.Semantics
	ctx.muUnlock()
	if sem != nil {
		return sem.Apply(tree, ruleName, params)
	}
	return tree, false
}
