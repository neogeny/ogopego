// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// TreeSeq represents a sequence node whose items are merged when folding.
type TreeSeq struct {
	Items []any
}

func (TreeSeq) tree() {}
