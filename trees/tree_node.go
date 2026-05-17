// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package trees

type Node struct {
	TreeBase
	TypeName string
	Tree     Tree
}

func (*Node) tree()                         {}
func (r *Node) fold(gather *treeMerge) Tree { return r }
