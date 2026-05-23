// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package ogopego

// allows embedding known grammars
import _ "embed"

// TatSuGrammarJSON is the embedded JSON representation of the Tatsu grammar
// used by the toolchain for bootstrapping and tests.
//
//go:embed grammar/tatsu.json
var TatSuGrammarJSON []byte
