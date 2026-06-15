// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// Seq represents a sequence node whose items are merged when folding.
type Seq struct {
	Items []any
}

func (Seq) tree() {}
func (s *Seq) fold(gather *treeMerge) any {
	var out any = nil
	for _, item := range s.Items {
		out = merge(out, fold(gather, item))
	}
	return out
}
