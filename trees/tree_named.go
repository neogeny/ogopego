package trees

import asjson "github.com/neogeny/ogopego/json"

type Named struct {
	TreeBase
	Name  string
	Value Tree
}

func (*Named) tree() {}
func (n *Named) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insert(n.Name, val)
	return val
}
func (n *Named) PubMap() *asjson.OrderedMap { return n.PubMapOf(n) }
func (n *Named) AsJSON() any {
	out := newOM()
	out.Set(n.Name, n.Value.AsJSON())
	return out
}
func (n *Named) AsJSONStr() string { return treeJSONStr(n.AsJSON()) }

type NamedAsList struct {
	TreeBase
	Name  string
	Value Tree
}

func (*NamedAsList) tree() {}
func (n *NamedAsList) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insertAsList(n.Name, val)
	return val
}
func (n *NamedAsList) PubMap() *asjson.OrderedMap { return n.PubMapOf(n) }
func (n *NamedAsList) AsJSON() any {
	out := newOM()
	out.Set(n.Name, n.Value.AsJSON())
	return out
}
func (n *NamedAsList) AsJSONStr() string { return treeJSONStr(n.AsJSON()) }

type Override struct {
	TreeBase
	Value Tree
}

func (*Override) tree() {}
func (o *Override) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendTree(gather.Root, val)
	return val
}
func (o *Override) PubMap() *asjson.OrderedMap { return o.PubMapOf(o) }
func (o *Override) AsJSON() any                { return o.Value.AsJSON() }
func (o *Override) AsJSONStr() string          { return treeJSONStr(o.AsJSON()) }

type OverrideAsList struct {
	TreeBase
	Value Tree
}

func (*OverrideAsList) tree() {}
func (o *OverrideAsList) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendAsList(gather.Root, val)
	return val
}
func (o *OverrideAsList) PubMap() *asjson.OrderedMap { return o.PubMapOf(o) }
func (o *OverrideAsList) AsJSON() any                { return o.Value.AsJSON() }
func (o *OverrideAsList) AsJSONStr() string          { return treeJSONStr(o.AsJSON()) }
