// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// Seq represents a sequence node whose items are merged when folding.
type Seq struct {
	TreeBase
	Items []Tree
}

func (Seq) tree() {}
func (s *Seq) fold(gather *treeMerge) Tree {
	var out Tree = &Nil{}
	for _, item := range s.Items {
		out = merge(out, item.fold(gather))
	}
	return out
}

// Array represents a closed list node produced after folding sequences.
type Array struct {
	TreeBase
	Items []Tree
}

func (Array) tree() {}
func (l *Array) fold(gather *treeMerge) Tree {
	items := make([]Tree, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.fold(gather)
	}
	return &Array{Items: items}
}
