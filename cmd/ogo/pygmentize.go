package main

import (
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

func Pygmentize(content, lexer string, useColor bool) string {
	if !useColor {
		return content
	}
	var buf strings.Builder
	err := quick.Highlight(&buf, content, lexer, "terminal256", "github-dark")
	if err != nil {
		return content
	}
	return buf.String()
}
