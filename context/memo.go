// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package context

import "github.com/neogeny/ogopego/trees"

type MemoKey struct {
	Mark    int
	Name    string
	CanMemo bool
}

type Memo struct {
	Tree trees.Tree
	Mark int
}
