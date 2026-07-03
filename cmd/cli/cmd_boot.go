package cli

import (
	"fmt"
	"os"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/util"
)

func bootCmd(cli CLIConfig) (string, []tOutputItem) {
	var outputs []tOutputItem
	gram, err := api.BootGrammar()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading boot grammar:", err)
		os.Exit(1)
	}
	var payload string
	var lang string
	switch {
	case cli.Boot.Json:
		payload = peg.ModelToJSONStr(gram)
		lang = "json"
	case cli.Boot.Pretty:
		payload = gram.PrettyPrint()
		lang = "ebnf"
	case cli.Boot.Railroads:
		payload = gram.Railroads()
		lang = "apl"
	case cli.Boot.Model:
		payload = util.Repr(gram)
		lang = "go"
	default:
		payload = gram.PrettyPrint()
		lang = "ebnf"
	}
	outputs = append(outputs, tOutputItem{Path: "boot", Output: payload})
	return lang, outputs
}
