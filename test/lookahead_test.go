package test

import (
	"testing"
)

func TestSkipTo(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Test
		start = 'x' ab $
		ab = 'a' 'b' | -> 'b'
	`), nil)
	AssertJSONStr(t, g, "x yb", `["x", "b"]`)
}
