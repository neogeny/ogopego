package test

import (
	"testing"
)

func TestNamedCapture(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := name='hello'
	`, nil)
	AssertJSONStr(t, g, "hello", `{"name": "hello"}`)
}

func TestOverrideSingleton(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start: ='hello'
	`, nil)
	AssertJSONStr(t, g, "hello", `"hello"`)
}
