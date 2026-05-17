// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package input

// Patterns defines regular expression patterns for whitespace, comments, and EOL comments.
type Patterns struct {
	Whitespace  string // Whitespace is the regular expression pattern for whitespace.
	Comments    string // Comments is the regular expression pattern for multi-line comments.
	EOLComments string // EOLComments is the regular expression pattern for end-of-line comments.
}
