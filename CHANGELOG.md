# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



[Unreleased]: https://github.com/neogeny/ogopego/compare/v0.1.8...HEAD
[0.1.8]: https://github.com/neogeny/ogopego/compare/v0.1.6...v0.1.8
[0.1.6]: https://github.com/neogeny/ogopego/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/neogeny/ogopego/compare/v0.1.2...v0.1.5
[0.1.2]: https://github.com/neogeny/ogopego/compare/v0.1.0...v0.1.2

## [Unreleased]

## [0.1.8] 2026-05-30 Python 

### Added

* A Python package with an _out-of-process_ integration of **OGoPEGo** was created and published to [PyPi]. **OGoPEGo** outperforms **TatSu** (Python) by 1.6x and **TieXiu** (Rust) by 2.5x in the benchmarks.

[PyPi]: https://pypi.org/project/ogopego/

### Changed

- Renamed `DisasterReport` to `ParseFailure` and removed the lightweight `Nope` error.

### In Progress
- The implementation of concurrent `Choice` options is in, but.. it doesn't work!


## [0.1.6] 2026-05-27 Debug

### Changed
- Solved more issues with `NameGuard`, which has complex semantics in legacy **TatSu**. 
    * `NameGuard` is `false` by default .
    * Setting a non-empty pattern for `whitespace` enables `NameGuard` in the asumption that the grammar wants token delimitation. 
    * Setting an _empty_ `whitespace` pattern leaves `NameGuard` unchanged.
    * A `Cursor` has a default `whitespace` of `(?m)\s+` for the benefit of new users, which are surprized if there is no-tokenization, but the default _does not_ enable `NameGuard`.
    * `NameGuard` may be set by the `@@nameguard` grammar directive or through `Cfg.nameguard`. The explicit value is nonored always.
  
- Forked and patched `github.com/dlclark/regexp2` because _~50%_ of both CPU and RAM were being consumed in its allocation of a `rune` slice for the input text for matching. The issue is reported [here](https://github.com/dlclark/regexp2/issues/103), and the pull request with the patch is [here](https://github.com/dlclark/regexp2/pull/104).

- Re-implemented `Cursor.MatchToken()` using only string operations, without regular expression tricks. Peformance improved considerably with the change.


## [0.1.5] 2026-05-26 Optimize

### Changed
- Implement `NameGuard` semantics with string and rune comparisons, replacing the previous implementation that used regular expressions. `NameGuard` avoids matching a token when it is not a complete word in the input, like not matching`"new"` when the input at the cursor is `newVar...`. 

  `NameGuard`can be activated with the `@nameguard:: true` grammar directive, or with the `Cfg.nameguard` configuration option. It is activated by default for grammars that define patterns for whitespace or comments.
- Use `BoundedMap` for the `Memo` cache. Parsing is faster with with a smaller cache that speeds up `Memo` lookups. The cache capacity is calculared using the heuristic `Cfg.PerLineMemos * Cursor.LineCount`. 
- The `Memo` cache is pruned when a `Cut`expression is parsed. Entries with marks lower than that of the previous cut are removed if they are not failure (`Tree.Bottom`) markers.

### Fixed
- Verified non-local/third-party builds through GitHub workflow.
- To be friendly with the ecosystem, skip the use of a `vendor` directory. A `_vendor` directory remains in the repo to guarantee build stability if a dependency becomes unreliable.

## [0.1.2] 2026-05-23 - Release


### Added
- Initial public release, feature-complete.
