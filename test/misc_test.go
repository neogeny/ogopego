package test

import (
	"testing"
)

func TestMapping(t *testing.T) {
	g := Compile(t, Dedent(`
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = key:key value:value
		key = /\w+/
		value = /\w+/
	`), nil)
	// Named items in a sequence are auto-folded into a single map.
	AssertJSONStr(t, g, "foo bar", `{"key": "foo", "value": "bar"}`)
}
