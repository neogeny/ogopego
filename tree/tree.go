package tree

type Tree interface {
	tree()
	fold(gather *treeMerge) Tree
}

type treeMerge struct {
	Root Tree
	Map  map[string]Tree
}

func Fold(tree Tree) Tree {
	if tree == nil {
		return &Nil{}
	}
	g := &treeMerge{Map: make(map[string]Tree)}
	result := tree.fold(g)
	return finish(g, result)
}

func finish(g *treeMerge, base Tree) Tree {
	switch g.Root.(type) {
	case *Nil, nil:
	default:
		return closed(g.Root)
	}
	if len(g.Map) > 0 {
		return &MapNode{Entries: g.Map}
	}
	return closed(base)
}

func closed(t Tree) Tree {
	if _, ok := t.(*Seq); ok {
		return &List{Items: t.(*Seq).Items}
	}
	return t
}

func merge(a, b Tree) Tree {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		sa, aOk := a.(*Seq)
		sb, bOk := b.(*Seq)
		switch {
		case aOk && bOk:
			items := make([]Tree, len(sa.Items)+len(sb.Items))
			copy(items, sa.Items)
			copy(items[len(sa.Items):], sb.Items)
			return &Seq{Items: items}
		case aOk:
			items := make([]Tree, len(sa.Items)+1)
			copy(items, sa.Items)
			items[len(sa.Items)] = b
			return &Seq{Items: items}
		case bOk:
			items := make([]Tree, 1+len(sb.Items))
			items[0] = a
			copy(items[1:], sb.Items)
			return &Seq{Items: items}
		default:
			return &Seq{Items: []Tree{a, b}}
		}
	}
}

func appendTree(a, b Tree) Tree {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		if s, ok := a.(*Seq); ok {
			items := make([]Tree, len(s.Items)+1)
			copy(items, s.Items)
			items[len(s.Items)] = b
			return &Seq{Items: items}
		}
		return &Seq{Items: []Tree{a, b}}
	}
}

func appendAsList(a, b Tree) Tree {
	if isNil(a) {
		return &Seq{Items: []Tree{b}}
	}
	if s, ok := a.(*Seq); ok {
		items := make([]Tree, len(s.Items)+1)
		copy(items, s.Items)
		items[len(s.Items)] = b
		return &Seq{Items: items}
	}
	return &Seq{Items: []Tree{a, b}}
}

func (m *treeMerge) insert(key string, val Tree) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = val
	} else {
		m.Map[key] = appendTree(existing, val)
	}
}

func (m *treeMerge) insertAsList(key string, val Tree) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = &Seq{Items: []Tree{val}}
	} else {
		m.Map[key] = appendAsList(existing, val)
	}
}

func isNil(t Tree) bool {
	_, ok := t.(*Nil)
	return ok || t == nil
}
