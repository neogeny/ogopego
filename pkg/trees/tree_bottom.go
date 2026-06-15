// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

type typeBOTTOM interface {
	Tree
	isBottom() bool
}

var BOTTOM typeBOTTOM = &typeBottomTree{}

// typeBottomTree is an internal sentinel node.
type typeBottomTree struct {
}

func (typeBottomTree) tree() {}

func (typeBottomTree) isBottom() bool {
	return true
}
