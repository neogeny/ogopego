# ogopego Style Guide

## Principles

The specific rules below derive from these higher-order principles. When in doubt, prefer the principle to the rule.

### Separation of Concerns

Every package, module, function, and type has exactly one responsibility. A package that reads files should not parse them. A function that builds a data structure should not format it for display. Transport types (`tPayload`, `tOutputItem`) exist specifically to cross a boundary between two concerns — they are pure data with no behavior.

**Violation (negative example):** Inlining output formatting logic inside a parser function, or having a single function that both processes results and writes them to the terminal.

The CLI layer (`cmd/cli/`) is the only place that knows about progress bars, color, and terminal width. The `pkg/parproc` package knows about concurrency and error collection but nothing about what it is processing or why.

### Make Invalid States Unrepresentable

A type's fields should have no (nil, zero, empty-string, or impossible-combination) states that a caller must avoid at runtime. Favor sum types (separate concrete types) over a single struct with an always-check-this-field approach.

**Example:** `ParProc.Result` has `Outcome` and `Err` fields. A caller checks `result.Err != nil` — but that's a leak. Better would be a discriminated type, but Go lacks sum types. The compromise: document clearly that exactly one of `Err` or `Outcome` is non-zero, and provide `func (r Result[P, R]) OK() bool { return r.Err == nil }`.

**Negative example:** A single `error` field that might be nil, plus a `success bool` field that must agree with it, plus a `reason string` field that's only valid when `success == false`.

### Explicit Over Implicit

No `init()` functions. No magic global state. No "and then this side effect happens" without a clear call site. Constructors (`New*`) do the work. Dependencies are wired explicitly.

### Favor Composition

Small, focused types that combine by holding references to each other — not by embedding or inheritance. Embedding is reserved for interface satisfaction (e.g., `json.Marshaler` delegation). Shared mutable state is a named pointer field (`mu *sync.Mutex`), not magic sharing via embedding.

### Conventions Before Configuration

Standard patterns reduce decision fatigue. Use the established file layout, naming scheme, and error-handling patterns. Only introduce configuration knobs when a convention demonstrably doesn't fit.

### Responsibility-Driven Placement

A symbol lives in the package that owns its primary concern. If you're adding a field, ask: "what package should own this?" If the answer is ambiguous, the design may need refactoring.

---

## Naming

### Types

- **Domain types** are bare nouns: `CoreCtx`, `MemoCache`, `Grammar`, `Cfg`, `ParseFailure`, `TokenStack`.
- **Transport types** (pure data crossing a boundary — between packages, goroutines, or pipeline stages) get a `t` prefix: `tPayload`, `tOutputItem`. These have no methods, only exported fields.
- **Interfaces** are short (often 1–3 methods): `Ctx` (not `Context`), `Heart` (not `Heartbeater`). Defined at the consumer site, in the package that uses them.
- **Acronyms** are all-caps: `EOF`, `EOL`, `JSON`, `HTML`, `ID`.
- **Type aliases** (`type X = Y`) decouple packages without wrapping: `type OrderedMap = orderedmap.OrderedMap[string, any]`.

### Functions and Methods

- **Constructor** `New*` returns a pointer: `NewCtx(...) *CoreCtx`.
- **Method receivers**: value receiver when the type is a lightweight handle or the method doesn't mutate. Pointer receiver once any method needs mutation.
- **Sentinel errors** prefixed with `Err` (`ErrPruningDisabled`). Errors are returned, never panicked.
- **`_` prefix** for private helpers that exist purely to decompose a public API function: `_newMemoKey`.

### Variables

- **Short in short scopes**: `ctx`, `cfg`, `lc`, `fp`, `sem`, `mu`, `wg`.
- **Descriptive in wider scopes**: `loadCfg`, `parseTask`, `payloads`, `maxWorkers`.
- **`var` for zero-value**, `:=` for non-zero initialization.

### Files

- One file per primary type group: `ctx_core.go` → `CoreCtx`, `memo.go` → `MemoCache`/`MemoKey`/`Memo`.
- `model_*.go` for PEG grammar node types (`model_choice.go`, `model_rule.go`).
- Test files: `foo_test.go` next to `foo.go`.

### Export Boundary

Everything is unexported by default. Export only when a symbol must cross a package boundary. The `api/` package is the public surface; `pkg/` exports what `api/` needs, nothing more.

---

## Package Layout

```
api/          — Public surface (LoadGrammar, ParseInput, Compile, BootGrammar)
pkg/          — Internal implementation
  ├── config/       — Configuration types
  ├── context/      — Parser runtime context
  ├── input/        — Cursor, input handling
  ├── asjson/       — AST-to-JSON conversion
  ├── peg/          — Grammar model, compilation, parsing
  ├── trees/        — Tree utilities
  ├── util/         — General-purpose utilities
  │   ├── container/
  │   ├── heartbeat/
  │   ├── newlines/
  │   └── pyre/
  ├── parproc/      — Generic parallel processing
  └── tool/         — Code generation
cmd/          — CLI entry points (thin)
test/         — Integration tests
```

Each package holds one concern. If a package accumulates unrelated responsibilities, split it.

---

## Struct Design

- **Zero values are valid where possible**; otherwise `New*` enforces invariants.
- **No `init()`** — constructors do the work.

### Pointer vs. Value Fields

- Shared mutable state: `mu *sync.Mutex`, `memoCache *MemoCache` — named pointers, not embedded.
- Immutable or locally-owned data: value fields.
- **Clone pattern**: share pointer fields explicitly, copy value fields:

```go
func (ctx *CoreCtx) Clone() Ctx {
    return &CoreCtx{
        cursor:    ctx.cursor.Clone(),
        mu:        ctx.mu,           // shared
        memoCache: ctx.memoCache,    // shared
        startedAt: ctx.startedAt,    // copied
    }
}
```

---

## Generics

- **Single-letter type params** describe the param's role: `P` (payload), `R` (result), `K`/`V` (key/value).
- **Generic utilities, concrete domain models** — `ParProc[P, R any]`, `OrderedMapEntries[K, V]`, `Chunks[T]` eliminate duplication in infrastructure. Domain types (`CoreCtx`, `Grammar`, `MemoCache`) are concrete.
- **No generic structs in domain code.**

---

## Concurrency

- **Thread-safety via explicit locking**, gated on a config flag:

```go
func (ctx *CoreCtx) muLock() {
    if ctx.cfg.Concurrency { ctx.mu.Lock() }
}
```

- **Bounded concurrency with a semaphore channel**: `sem := make(chan int, n)`.
- **Cancellation via `done` channel**: `close(done)` signals workers; each selects before proceeding.
- **Leak prevention**: `defer wg.Done()`, `defer func() { <-sem }()` always in scope.
- **Streaming results via `iter.Seq`**: buffered channel; consumer loop drives completion.
- **`sync.Mutex` is private** — never exposed on the type signature, never copied.

---

## Imports

Three groups, blank-line separated:

```go
import (
    "fmt"
    "os"
    "sync"

    "github.com/fatih/color"
    "github.com/vbauerster/mpb/v8"

    "github.com/neogeny/ogopego/api"
    "github.com/neogeny/ogopego/pkg/config"
)
```

- **Aliases** only to resolve ambiguity: `mpb "github.com/vbauerster/mpb/v8"`.
- **No `.` imports** except test data.
- **No blank imports** — if required, document with a comment.

---

## Error Handling

- **`errors.AsType[*T]`** for type-asserting wrapped errors.
- **Custom error types** carry structured context: `ParseFailure` has `Inner`, `Memento`, `Location`.
- **Guard-clause style** — early return on error, happy path after:

```go
if err != nil {
    fp.Fail()
    return nil, err
}
fp.Success()
```

- **Failures as values** — `ctx.Failure(mark, error)` returns an `error`, never panics.
- **`_` for unused returns** — explicit: `_, _ = fmt.Fprintf(...)`.

---

## Comments

- **Godoc on every exported symbol**: single-line summary. Package doc on the `package` clause.
- **Inline comments** explain "why" not "what." No trailing noise.
- **No comments on unexported symbols** unless explaining invariants or non-obvious design decisions.

---

## Testing

- **File placement**: `foo_test.go` in the same package (`package foo`) for white-box tests; external (`package foo_test`) for black-box.
- **Table-driven**:

```go
tests := []struct {
    name  string
    input string
    want  int
}{
    {name: "empty", input: "", want: 0},
    ...
}
```

- **`testing.TB` parameter** on test helpers: `func Compile(t testing.TB, ...)`.
- **`t.Helper()`** on every test utility.
- **`assert` package** — `github.com/alecthomas/assert/v2`:

```go
assert.Equal(t, want, got)
assert.NoError(t, err, "parse %q", text)
assert.Error(t, err)
```

- **Test data inline** — small grammars and inputs as backtick literals or `Dedent(...)`. No fixture files for unit tests.

---

## Coding Style

- **Boolean chains use `switch`**: `switch { case cond1: ... case cond2: ... }`.
- **Type switches** for `any` dispatch.
- **`var` blocks** for grouped declarations.
- **Explicit struct keys** in composite literals — positional only for trivial (≤2 field) structs.
- **Constants over magic numbers** — `maxRecursionDepth = 64`.
- **Nil slices are empty** — `for range slice` is safe on nil.
- **No naked returns** — always specify return values.
- **`_` for unused params** in interface implementations: `func Beat(mark, _ int)`.

---

## Build System

- **Targets** via `Justfile`: `just build`, `just lint`, `just vet`, `just test`.
- **Pipeline**: `gofmt -s` → `go vet` → `golangci-lint` → `go test` → `go build`.
- **Toolchain**: `go 1.26.3` — use idioms: `iter.Seq`, `errors.AsType`, `unique.Handle`, `slices`, `maps`.
- **`golangci-lint`** with `--exclude-dirs ./tmp`.
- **`gotestsum`** for test output.
- **Vendoring** (`internal/_vendor`) only for release builds.

---

## Copyright

```
// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0
```

On every source file.
