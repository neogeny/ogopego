package trees

import asjson "github.com/neogeny/ogopego/json"

type Node struct {
	TreeBase
	TypeName string
	Tree     Tree
}

func (*Node) tree()                         {}
func (r *Node) fold(gather *treeMerge) Tree { return r }
func (r *Node) PubMap() *asjson.OrderedMap  { return r.PubMapOf(r) }
func (r *Node) AsJSON() any {
	child := r.Tree.AsJSON()
	if m, ok := child.(map[string]any); ok {
		if _, has := m["__class__"]; !has {
			out := make(map[string]any, len(m)+1)
			out["__class__"] = r.TypeName
			for k, v := range m {
				out[k] = v
			}
			return out
		}
	}
	return map[string]any{"__class__": r.TypeName, "ast": child}
}
func (r *Node) AsJSONStr() string { return treeJSONStr(r.AsJSON()) }
