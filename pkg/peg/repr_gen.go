package peg

import (
	"strings"

	"github.com/neogeny/ogopego/pkg/util"
)

func ParserRepr(g Grammar, pkg string) string {
	funcName := g.Name + "Parser"

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkg)
	buf.WriteString("\n\nimport peg \"github.com/neogeny/ogopego/pkg/peg\"\n\nfunc ")
	buf.WriteString(funcName)
	buf.WriteString("() peg.Grammar {\n\treturn ")

	repr := util.Repr(g)
	lines := strings.Split(repr, "\n")
	for i, line := range lines {
		if i > 0 {
			buf.WriteString("\n\t")
		}
		buf.WriteString(line)
	}

	buf.WriteString("\n}\n")
	return buf.String()
}
