// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// treeSeq represents a sequence node whose items are merged when folding.
type treeSeq struct {
	Items []any
}

func (treeSeq) isTree() {}
