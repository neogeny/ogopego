package context

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/heartbeat"
)

// CoreCtx is the concrete implementation of Ctx used by the parser runtime.
type CoreCtx struct {
	cfg            Cfg
	cursor         Cursor
	callStack      CallStack
	cutStack       []bool
	tracer         Tracer
	furthest       *DisasterReport
	keywords       map[string]struct{}
	memoCache      MemoCache
	recursionKey   MemoKey
	recursionDepth int
	lookaheadDepth int
	lastCutMark    int
	heartbeat      heartbeat.Heartbeat
	heartbeatTime  time.Time
}

// NewCtx creates a new CoreCtx backed by the provided Cursor and optional
// configuration. Use the returned value where a context implementing Ctx is
// required.
func NewCtx(cursor Cursor, cfg *Cfg) *CoreCtx {
	stackCapacity := 64
	cfgS := cfg.New()
	memoCapacity := max(
		stackCapacity,
		int(math.Round(cfgS.PerLineMemos*float64(cursor.LineCount()))),
	)
	ctx := CoreCtx{
		// this is how passed configuration gets injected
		cfg:       cfgS,
		cursor:    cursor,
		tracer:    NullTracer{},
		heartbeat: heartbeat.NullHeartbeat{},
		callStack: make(CallStack, 0, stackCapacity),
		cutStack:  make([]bool, 1, stackCapacity),
		memoCache: NewMemoMache(memoCapacity),
	}
	ctx.cursor.Configure(ctx.cfg)
	ctx.cfg = ctx.cfg.Override(cfg)
	return &ctx
}

func (ctx *CoreCtx) Cfg() Cfg { return ctx.cfg }

func (ctx *CoreCtx) Configure(cfg Cfg) {
	ctx.cfg = ctx.cfg.Override(&cfg)
	ctx.cursor.Configure(cfg)
	ctx.setKeywords(cfg.Keywords)

	if cfg.Trace {
		if cfg.Colorize {
			color.Output = os.Stderr
			color.NoColor = false
		}
		ctx.tracer = ConsoleTracer{}
	} else {
		ctx.tracer = NullTracer{}
	}
	// FIXME  should have been set by Override()
	if cfg.Heartbeat != nil {
		ctx.heartbeat = cfg.Heartbeat
	}
}

func (ctx *CoreCtx) SetTracer(tracer Tracer) {
	ctx.tracer = tracer
}

func (ctx *CoreCtx) Cursor() Cursor { return ctx.cursor }

func (ctx *CoreCtx) CallStack() CallStack {
	cs := make([]string, len(ctx.callStack))
	copy(cs, ctx.callStack)
	return cs
}

func (ctx *CoreCtx) Tracer() Tracer { return ctx.tracer }

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
	if time.Since(ctx.heartbeatTime) < 128*time.Millisecond {
		return
	}
	mark := ctx.Mark()
	total := len(ctx.cursor.AsStr())
	if total == 0 {
		return
	}
	ctx.heartbeat.Tick(mark, total)
	ctx.heartbeatTime = time.Now()
}

func (ctx *CoreCtx) pruneCache(cutPoint int) {
	if ctx.cfg.NoPruneMemosOnCut {
		return
	}
	PruneMemoCache(ctx.memoCache, cutPoint)
}

func (ctx *CoreCtx) Key(name string, canMemo bool) MemoKey {
	return MemoKey{Mark: ctx.Mark(), Name: name, CanMemo: canMemo}
}

func (ctx *CoreCtx) Memo(key MemoKey) (Memo, bool) {
	return ctx.memoCache.Get(key)
}

func (ctx *CoreCtx) Memoize(key MemoKey, tree trees.Tree, mark int) {
	if !key.CanMemo {
		return
	}
	ctx.memoCache.Set(key, Memo{Tree: tree, Mark: mark})
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
	_, ok := ctx.keywords[name]
	return ok
}

func (ctx *CoreCtx) setKeywords(keywords []string) {
	sort.Strings(keywords)
	set := make(map[string]struct{}, len(keywords))
	for _, kw := range keywords {
		set[kw] = struct{}{}
	}
	ctx.keywords = set
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
	ctx.callStack = append(ctx.callStack, name)
}

func (ctx *CoreCtx) Leave() {
	if len(ctx.callStack) > 0 {
		ctx.callStack = ctx.callStack[:len(ctx.callStack)-1]
	}
}

func (ctx *CoreCtx) ParseEOF() bool {
	ctx.Enter("__eof__")
	ctx.NextToken()
	result := ctx.cursor.AtEnd()
	ctx.Leave()
	return result
}

func (ctx *CoreCtx) Failure(start int, source error) *Nope {
	loc := Location{}
	_, loc.File, loc.Line, _ = runtime.Caller(2)

	ctx.Reset(start)
	nope := &Nope{
		location: loc,
	}
	if furthest := ctx.FurthestFailure(); furthest != nil &&
		furthest.Start() >= start {
		return nope
	}
	msg := source.Error()
	dis := &DisasterReport{
		location: loc,
		Inner:    source,
		CutSeen:  ctx.IsCutSeen(),
		Memento: NewMemento(
			start,
			msg,
			ctx.cursor,
			ctx.CallStack(),
		),
	}
	ctx.SetFurthestFailure(dis)
	return nope
}

func (ctx *CoreCtx) FurthestFailure() *DisasterReport { return ctx.furthest }

func (ctx *CoreCtx) SetFurthestFailure(dis *DisasterReport) { ctx.furthest = dis }

func (ctx *CoreCtx) MatchPattern(pattern string) (string, error) {
	mark := ctx.Mark()
	m, ok := ctx.cursor.MatchPattern(pattern)
	if !ok {
		return "", ctx.Failure(mark, fmt.Errorf("expected pattern %q", pattern))
	}
	return m, nil
}

func (ctx *CoreCtx) Void() {
	ctx.NextToken()
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

func (ctx *CoreCtx) Constant(literal any) (trees.Tree, error) {
	switch v := literal.(type) {
	case string:
		return &trees.Text{Value: v}, nil
	case float64:
		return &trees.Number{Value: v}, nil
	case bool:
		return &trees.Bool{Value: v}, nil
	case nil:
		return trees.NIL, nil
	case int:
		return &trees.Number{Value: float64(v)}, nil
	default:
		return &trees.Text{Value: fmt.Sprintf("%v", v)}, nil
	}
}

func (ctx *CoreCtx) Cut() {
	ctx.cutStack[len(ctx.cutStack)-1] = true
	ctx.Tracer().TraceCut(ctx)
	if ctx.lookaheadDepth == 0 {
		mark := ctx.Mark()
		if mark > ctx.lastCutMark {
			ctx.pruneCache(ctx.lastCutMark)
			ctx.lastCutMark = mark
		}
	}
}

func (ctx *CoreCtx) SetCut() {
	ctx.cutStack[len(ctx.cutStack)-1] = true
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

func (ctx *CoreCtx) ApplySemantics(node trees.Tree, ruleName string, params []string) (trees.Tree, bool) {
	if ctx.cfg.Semantics != nil {
		return ctx.cfg.Semantics(node, ruleName, params)
	}
	return node, false
}
