# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

[Unreleased]: https://github.com/neogeny/ogopego/compare/v0.1.2...HEAD

### Added

### Fixed

## [0.1.3] - 2026-05-25

[0.1.3]: https://github.com/neogeny/ogopego/compare/v0.1.2...v0.1.3

### Changed
- Implement `NameGuard` semantics with string and rune comparisons, replacing the previous implementation that used regular expressions. `NameGuard` avoids matching a token when it is not a complete word in the input, like not matching`"new"` when the input at the cursor is `newVar...`. 

  `NameGuard`can be activated with the `@nameguard:: true` grammar directive, or with the `Cfg.nameguard` configuration option. It is activated by default for grammars that define patterns for whitespace or comments.
- Use `BoundedMap` for the `Memo` cache. Parsing is faster with with a smaller cache that speeds up `Memo` lookups. The cache capacity is calculared using the heuristic `Cfg.PerLineMemos * Cursor.LineCount`. 
- The `Memo` cache is pruned when a `Cut`expression is parsed. Entries with marks lower than that of the previous cut are removed if they are not failure (`Tree.Bottom`) markers.

## [0.1.2] - 2026-05-23

[0.1.2]: https://github.com/neogeny/ogopego/compare/v0.1.0...v0.1.2

### Added
- Initial public release, feature-complete.
