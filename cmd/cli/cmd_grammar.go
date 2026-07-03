package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/tool"
	"github.com/neogeny/ogopego/pkg/util"
	"github.com/vbauerster/mpb/v8"
)

func grammarCmd(cli CLIConfig, cliCfg *config.Cfg) (string, []tOutputItem) {
	var outputs []tOutputItem
	var p *mpb.Progress
	if !cli.Quiet {
		hideCursor()
		p = mpb.New(mpb.WithOutput(os.Stderr))
	}
	fileName := filepath.Base(cli.Grammar.Grammar)
	fp := NewFileProgress(p, fileName)

	data, err := os.ReadFile(cli.Grammar.Grammar)
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nerror reading grammar:", err)
		os.Exit(1)
	}
	fp.SetLength(len(data))

	loadCfg := *cliCfg
	loadCfg.Heart = fp.Heartbeat()
	gram, err := api.LoadGrammar(cli.Grammar.Grammar, &loadCfg)
	if err != nil {
		fp.Fail()
		if p != nil {
			p.Wait()
			showCursor()
		}
		fmt.Fprintln(os.Stderr, "\nerror:", err)
		os.Exit(1)
	}
	fp.Success()
	if p != nil {
		p.Wait()
		showCursor()
	}
	var output string
	var lang string
	switch {
	case cli.Grammar.Json:
		output = peg.ModelToJSONStr(gram)
		lang = "json"
	case cli.Grammar.Pretty:
		output = gram.PrettyPrint()
		lang = "ebnf"
	case cli.Grammar.Railroads:
		output = gram.Railroads()
		lang = "apl"
	case cli.Grammar.Model:
		output = util.Repr(gram)
		lang = "go"
	case cli.Grammar.Parser != "":
		output = peg.ParserRepr(*gram, cli.Grammar.Parser)
		lang = "go"
	case cli.Grammar.ModelGen != "":
		output = tool.GenerateGrammarModel(*gram, cli.Grammar.ModelGen)
		lang = "go"
	default:
		output = gram.PrettyPrint()
		lang = "ebnf"
	}
	outputs = append(outputs, tOutputItem{Path: cli.Grammar.Grammar, Output: output})
	return lang, outputs
}
