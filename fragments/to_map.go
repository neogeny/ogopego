// The Mirror: Converts OrderedMap tree to a standard map tree
func ToMap(v any) any {
	switch val := v.(type) {
	case *orderedmap.OrderedMap:
		m := make(map[string]any)
		for _, k := range val.Keys() {
			item, _ := val.Get(k)
			m[k] = ToMap(item)
		}
		return m
	case []any:
		s := make([]any, len(val))
		for i, item := range val {
			s[i] = ToMap(item)
		}
		return s
	default:
		return v
	}
}
