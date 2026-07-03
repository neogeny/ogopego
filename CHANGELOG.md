#.Tree Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

[TatSu]: https://github.com/neogeny/TatSu
[TieXiu]: https://github.com/neogeny/tiexiu
[OGoPEGo]: https://github.com/neogeny/ogopego
[TSemekwes]: https://github.com/neogeny/tsemekwes
[PyPi]: https://pypi.org/project/ogopego/


[Unreleased]: https://github.com/neogeny/ogopego/compare/v0.1.13...HEAD
[v0.1.13]: https://github.com/neogeny/ogopego/compare/v0.1.11...v0.1.13
[v0.1.12]: https://github.com/neogeny/ogopego/compare/v0.1.11...v0.1.12
[v0.1.11]: https://github.com/neogeny/ogopego/compare/v0.1.9...v0.1.11
[v0.1.9]: https://github.com/neogeny/ogopego/compare/v0.1.6...v0.1.9
[v0.1.6]: https://github.com/neogeny/ogopego/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/neogeny/ogopego/compare/v0.1.2...v0.1.5
[v0.1.2]: https://github.com/neogeny/ogopego/compare/v0.1.0...v0.1.2

## [Unreleased]

### Added

- Generic `ParProc[P, R any]` parallel processor in `pkg/parproc/` with bounded concurrency, cancellation, and per-task timing.
- `STYLE.md` — project-wide coding conventions captured from the codebase.

### Changed

- Refactored `cmd_run` to use a two-phase pipeline (sequential read, parallel parse) with `ParProc`.
- Eliminated `CoreCtxHeavy` — its fields are now inlined into `CoreCtx` behind explicit pointers for shared mutable state.
- Removed the experimental `sync.Mutex` from `CoreCtx`. Each `Clone()` now deep-copies mutable state, making every context fully isolated and safe for concurrent use.
- Revised progress-bar cosmetics (wider file bars, shorter labels, color tweaks) and added `progress.Write` for goroutine-safe output.
- Heartbeat interval now computed in beats-per-second (BPS).

#### Optimization

- The latest round of hunting for bottlenecks has **OGoPEGo** generating the fastest parsers in its family. The **OGoPEGo** model-based parser for Java is 1.18x faster than the one produced by [TatSu] with procedural code generation, the previous fastest. Without the out-of-process overhead when executed as a Python package, the Go parser is 4x faster than the best Python parser.
- The latest heuristics for optimizing grammar models were ported from [TatSu].

## [v0.1.13] 2026-06-14 New syntax + optimized grammar models

### Added

- New `@name`, `@int`, `@uint`, `@float`, `@bool` meta-expressions for typed token matching (names, signed/unsigned ints, floats, bools)
- `Cursor.MatchName`/`MatchInt`/`MatchUInt`/`MatchFloat`/`MatchBool` methods on `RuneCursor` and `StrCursor`
- `Ctx` interface methods that delegate to cursor matching
- Five concrete model types (`NameMeta`, `IntMeta`, `UIntMeta`, `FloatMeta`, `BoolMeta`) with `MetaExp` base
- Full compilation pipeline: EBNF tree compile, JSON import/export, analysis passes, display (pretty-print & railroads)
- 40 cursor-level tests ported from TatSu
- Boot grammar usage of `@name` in rule definitions

### Changed

- Meta expression dispatch updated from single `*MetaExp` type case to five individual type cases across the compiler pipeline
- All TatSu `@` meta features brought to parity with TatSu v5.21.1

### Changed

- Refactored grammar semantics from function type (`SemanticsFunc`) to `GrammarSemantics` interface in `pkg/config`; replaced `GrammarParserSemantics` with `EBNFGrammarSemantics` struct adding meta-expression node production (`NameMeta`, `IntMeta`, `UIntMeta`, `FloatMeta`, `BoolMeta`) during boot grammar compilation of `@` expressions

## [v0.1.12] 2026-06-06 Optimize

### Changed

* `Rule.Optimized()` recursively optimizes the expression tree, unwraps single-element `Sequence` containers, and inlines alias rules whose body is a single `Call`. `optimizeExpr` always clones nodes so the original grammar's expression tree remains immutable after `Initialize()` links them using the optimized `Rule`. `Model.Optimized()` is value-based so the original `Grammar` model remains unchanged.

## [v0.1.11] 2026-05-31 Optimize

### Changed

* Define `RuneCursor` and refactor `pyre.Pattern` and friends to take advantage of `rune` matching in [dlclark/regexp2/v2]. Performances improved considerably.

[dlclark/regexp2/v2]:https://pkg.go.dev/github.com/dlclark/regexp2/v2


* Converted tests to use [alecthomas/assert/v2] for better readability and maintainability. The new assertion library provides a more expressive and concise way to write tests, improving the overall quality of the test suite.

[alecthomas/assert/v2]: https://pkg.go.dev/github.com/alecthomas/assert/v2

* Renamed `trees.List` to `trees.Array` for matching semantics and closeness to Go and JSON.

* Added a progress bar to the `grammar` sub-command of the CLI.

## [v0.1.9] 2026-05-30 Python

### Added

* A Python package with an _out-of-process_ integration of **[OGoPEGo]** was created and published to [PyPi]. **OGoPEGo** outperforms **TatSu** (Python) by 1.6x and **TieXiu** (Rust) by 2.5x in the benchmarks.

### Changed

- Renamed `DisasterReport` to `ParseFailure` and removed the lightweight `Nope` error.

### In Progress
- The implementation of concurrent `Choice` options is in, but.. it doesn't work!


## [v0.1.6] 2026-05-27 Debug

### Changed
- Solved more issues with `NameGuard`, which has complex semantics in legacy **TatSu**.
    * `NameGuard` is `false` by default .
    * Setting a non-empty pattern for `whitespace` enables `NameGuard` in the asumption that the grammar wants token delimitation.
    * Setting an _empty_ `whitespace` pattern leaves `NameGuard` unchanged.
    * A `Cursor` has a default `whitespace` of `(?m)\s+` for the benefit of new users, which are surprized if there is no-tokenization, but the default _does not_ enable `NameGuard`.
    * `NameGuard` may be set by the `@@nameguard` grammar directive or through `Cfg.nameguard`. The explicit value is nonored always.

- Forked and patched `github.com/dlclark/regexp2` because _~50%_ of both CPU and RAM were being consumed in its allocation of a `rune` slice for the input text for matching. The issue is reported [here](https://github.com/dlclark/regexp2/issues/103), and the pull request with the patch is [here](https://github.com/dlclark/regexp2/pull/104).

- Re-implemented `Cursor.MatchToken()` using only string operations, without regular expression tricks. Peformance improved considerably with the change.


## [v0.1.5] 2026-05-26 Optimize

### Changed
- Implement `NameGuard` semantics with string and rune comparisons, replacing the previous implementation that used regular expressions. `NameGuard` avoids matching a token when it is not a complete word in the input, like not matching`"new"` when the input at the cursor is `newVar...`.

  `NameGuard`can be activated with the `@nameguard:: true` grammar directive, or with the `Cfg.nameguard` configuration option. It is activated by default for grammars that define patterns for whitespace or comments.
- Use `BoundedMap` for the `Memo` cache. Parsing is faster with a smaller cache that speeds up `Memo` lookups. The cache capacity is calculared using the heuristic `Cfg.PerLineMemos * Cursor.LineCount`.
- The `Memo` cache is pruned when a `Cut`expression is parsed. Entries with marks lower than that of the previous cut are removed if they are not failure (`Tree.BOTTOM`) markers.

### Fixed
- Verified non-local/third-party builds through GitHub workflow.
- To be friendly with the ecosystem, skip the use of a `vendor` directory. An `internal/_vendor` directory remains in the repo to guarantee build stability if a dependency becomes unreliable.

## [v0.1.2] 2026-05-23 Release


### Added
- Initial public release, feature-complete.
