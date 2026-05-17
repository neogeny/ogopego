// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package ogopego

import _ "embed"

// TatsuGrammarJSON is the embedded JSON representation of the Tatsu grammar
// used by the toolchain for bootstrapping and tests.
//
//go:embed grammar/tatsu.json
var TatsuGrammarJSON []byte
