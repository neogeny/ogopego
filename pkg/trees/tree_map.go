// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

// MapNode represents a keyed mapping of entries produced during folding.
type MapNode struct {
	Entries map[string]Tree
}

func (MapNode) tree()                          {}
func (m *MapNode) fold(gather *treeMerge) Tree { return m }
