package api

import (
	"sync"

	"github.com/neogeny/ogopego"
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/input"
	"github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/peg"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util/pyre"
)

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

func ParseGrammar(grammar string) (trees.Tree, error) {
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
	return boot.ParseTree(context.NewBaseCtx(c))
}

func ParseGrammarToJSON(grammar string) (any, error) {
	tree, err := ParseGrammar(grammar)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(tree), nil
}

func ParseGrammarToJSONString(grammar string) (string, error) {
	tree, err := ParseGrammar(grammar)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(tree), nil
}

func Compile(grammar string) (*peg.Grammar, error) {
	tree, err := ParseGrammar(grammar)
	if err != nil {
		return nil, err
	}
	return peg.Compile(tree)
}

func CompileToJSON(grammar string) (any, error) {
	g, err := Compile(grammar)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(g), nil
}

func CompileToJSONString(grammar string) (string, error) {
	g, err := Compile(grammar)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(g), nil
}

func LoadGrammarFromJSON(data []byte) (*peg.Grammar, error) {
	return json.ParseGrammar(data)
}

// func LoadGrammarFromJSONToJSON(data []byte) (any, error) { ... }

// func LoadTreeFromJSON(data []byte) (trees.Tree, error) { ... }
// func LoadTreeToJSON(data []byte) (any, error) { ... }

// func GrammarPretty(grammar string) (string, error) { ... }
//   needs Compile + pretty printing.

// func Parse(grammar, text string) (trees.Tree, error) { ... }
//   needs Compile.
// func ParseToJSON(grammar, text string) (any, error) { ... }
// func ParseToJSONString(grammar, text string) (string, error) { ... }

func ParseInput(parser *peg.Grammar, text string) (trees.Tree, error) {
	ctx := context.NewBaseCtx(input.NewStrCursor(text))
	return parser.ParseTree(ctx)
}

func ParseInputToJSON(parser *peg.Grammar, text string) (any, error) {
	tree, err := ParseInput(parser, text)
	if err != nil {
		return nil, err
	}
	return json.AsJSON(tree), nil
}

func ParseInputToJSONString(parser *peg.Grammar, text string) (string, error) {
	tree, err := ParseInput(parser, text)
	if err != nil {
		return "", err
	}
	return json.AsJSONs(tree), nil
}

// func BootGrammarToJSON() (any, error) { ... }
// func BootGrammarToJSONString() (string, error) { ... }
// func BootGrammarPretty() (string, error) { ... }
