package test

import (
	"testing"
)

func TestPositiveJoin(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@nameguard :: False
		@@grammar :: Test
		start := ','%{'x' 'y'}+
	`, nil)
	AssertJSONStr(t, g, "x y, x y", `[["x", "y"], ",", ["x", "y"]]`)
	AssertJSONStr(t, g, "x y x y", `[["x", "y"]]`)
	ParseFail(t, g, "y x")
}
