package trees

import asjson "github.com/neogeny/ogopego/json"

type MapNode struct {
	TreeBase
	Entries map[string]Tree
}

func (*MapNode) tree()                         {}
func (m *MapNode) fold(gather *treeMerge) Tree { return m }
func (m *MapNode) PubMap() *asjson.OrderedMap  { return m.PubMapOf(m) }
func (m *MapNode) AsJSON() any {
	out := make(map[string]any, len(m.Entries))
	for k, v := range m.Entries {
		out[k] = v.AsJSON()
	}
	return out
}
func (m *MapNode) AsJSONStr() string { return treeJSONStr(m.AsJSON()) }
