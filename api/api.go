package api

import (
	"sync"

	"github.com/neogeny/ogopego"
	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/peg"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

type Cfg = config.Cfg

var (
	bootOnce sync.Once
	bootGram *peg.Grammar
	bootErr  error
)

func bootGrammar() (*peg.Grammar, error) {
	bootOnce.Do(func() {
		bootGram, bootErr = json.LoadBootGrammar(ogopego.TatsuGrammarJSON)
	})
	return bootGram, bootErr
}

func BootGrammar() (*peg.Grammar, error) {
	return bootGrammar()
}

func ParseGrammar(grammar string, cfg *Cfg) (trees.Tree, error) {
	boot, err := bootGrammar()
	if err != nil {
		return nil, err
	}
	pat, err := pyre.Compile(`(?m)[ \t]+`)
	if err != nil {
		return nil, err
	}
	c := input.NewStrCursor(grammar)
	c.SetPatterns(&input.TokenizingPatterns{Wsp: pat})
	return boot.Parse(context.NewCtx(c, cfg), cfg)
}

func ParseGrammarToJSON(grammar string, cfg *Cfg) (any, error) {
	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(tree), nil
}

func ParseGrammarToJSONString(grammar string, cfg *Cfg) (string, error) {
	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(tree), nil
}

func Compile(grammar string, cfg *Cfg) (*peg.Grammar, error) {
	tree, err := ParseGrammar(grammar, cfg)
	if err != nil {
		return nil, err
	}
	return peg.Compile(tree)
}

func CompileToJSON(grammar string, cfg *Cfg) (any, error) {
	g, err := Compile(grammar, cfg)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(g), nil
}

func CompileToJSONString(grammar string, cfg *Cfg) (string, error) {
	g, err := Compile(grammar, cfg)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(g), nil
}

func ParseInput(parser *peg.Grammar, text string, cfg *Cfg) (trees.Tree, error) {
	ctx := context.NewCtx(input.NewStrCursor(text), cfg)
	return parser.Parse(ctx, cfg)
}

func ParseInputToJSON(parser *peg.Grammar, text string, cfg *Cfg) (any, error) {
	tree, err := ParseInput(parser, text, cfg)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(tree), nil
}

func ParseInputToJSONString(parser *peg.Grammar, text string, cfg *Cfg) (string, error) {
	tree, err := ParseInput(parser, text, cfg)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(tree), nil
}
func LoadGrammarFromJSON(data []byte) (*peg.Grammar, error) {
	return json.ParseGrammar(data)
}
