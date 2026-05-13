func Compare(actual any, expected any) bool {
	switch a := actual.(type) {
	case *orderedmap.OrderedMap:
		e, ok := expected.(*orderedmap.OrderedMap)
		if !ok || len(a.Keys()) != len(e.Keys()) {
			return false
		}
		for _, k := range a.Keys() {
			valA, _ := a.Get(k)
			valE, _ := e.Get(k)
			if !Compare(valA, valE) {
				return false
			}
		}
		return true

	case []any:
		e, ok := expected.([]any)
		if !ok || len(a) != len(e) {
			return false
		}
		for i := range a {
			if !Compare(a[i], e[i]) {
				return false
			}
		}
		return true

	default:
		// Handles string, int, bool, etc.
		return actual == expected
	}
}
