package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/peg"
)

func main() {
	ctx = kong.Parse(&CLI,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
		kong.Help(coloredHelp),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: false,
		}),
	)

	if err := validateExclusive("format"); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := ctx.Selected()

	switch CLI.Color {
	case "always":
		useColorOutput = true
		color.NoColor = false
	case "never":
		useColorOutput = false
		color.NoColor = true
	case "auto":
		if isTerminal() {
			useColorOutput = true
			color.NoColor = false
		} else {
			useColorOutput = false
			color.NoColor = true
		}
	}
	_ = useColorOutput

	var output string

	if cmd != nil {
		switch cmd.Name {
		case "run":
			err := fmt.Errorf("run command not fully wired yet")
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)

		case "boot":
			gram, err := api.BootGrammar()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error loading boot grammar:", err)
				os.Exit(1)
			}
			switch {
			case CLI.Boot.Json:
				output = json.AsJSONs(gram)
			case CLI.Boot.Pretty:
				output = gram.PrettyPrint()
			case CLI.Boot.Railroads:
				output = gram.Railroads()
			default:
				output = gram.PrettyPrint()
			}

		case "grammar":
			gram, err := loadGrammar(CLI.Grammar.Grammar)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			switch {
			case CLI.Grammar.Json:
				output = json.AsJSONs(gram)
			case CLI.Grammar.Pretty:
				output = gram.PrettyPrint()
			case CLI.Grammar.Railroads:
				output = gram.Railroads()
			default:
				output = gram.PrettyPrint()
			}
		}
	}

	if output != "" {
		if CLI.Output != "" {
			if err := os.WriteFile(CLI.Output, []byte(output), 0644); err != nil {
				fmt.Fprintln(os.Stderr, "error writing output:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(output)
		}
	}
}

func loadGrammar(path string) (*peg.Grammar, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return api.LoadGrammarFromJSON(data)
	default:
		return api.Compile(string(data), nil)
	}
}
