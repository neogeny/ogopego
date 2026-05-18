// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/peg"
	"github.com/neogeny/ogopego/util"
)

func main() {
	ctx = kong.Parse(&CLI,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
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
	cliCfg = &config.Cfg{
		Trace:    CLI.Trace,
		Colorize: useColorOutput,
	}

	var output string
	var lang string

	if cmd != nil {
		switch cmd.Name {
		case "run":
			prog := NewProgressUI(len(CLI.Run.Inputs))
			loader := prog.Loading("loading grammar")
			loadCfg := *cliCfg
			loadCfg.Heartbeat = loader.Heartbeat()
			gram, err := loadGrammar(CLI.Run.Grammar, &loadCfg)
			loader.Finish()
			if err != nil {
				fmt.Fprintln(os.Stderr, "\nerror:", err)
				os.Exit(1)
			}
			var errcount int
			for _, path := range CLI.Run.Inputs {
				name := filepath.Base(path)
				fp := prog.AddFile(name)

				data, err := os.ReadFile(path)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "\nerror reading %s: %v\n", path, err)
					errcount++
					prog.IncFiles()
					continue
				}
				fp.SetLength(len(data))

				fileCfg := *cliCfg
				fileCfg.Heartbeat = fp.Heartbeat()
				result, err := api.ParseInputToJSONString(gram, string(data), &fileCfg)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "\nerror parsing %s: %v\n", path, err)
					errcount++
					fp.Fail()
				} else {
					fp.Success()
					output += result + "\n"
				}
				prog.IncFiles()
			}
			prog.Finish()
			passed := len(CLI.Run.Inputs) - errcount
			_, _ = fmt.Fprintf(os.Stderr, "%s %s %s\n",
				color.New(color.FgWhite, color.Bold).Sprintf("Parsed %d files", len(CLI.Run.Inputs)),
				color.New(color.FgGreen, color.Bold).Sprintf("%d passed", passed),
				color.New(color.FgRed, color.Bold).Sprintf("%d errors", errcount),
			)
			lang = "json"

		case "boot":
			gram, err := api.BootGrammar()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error loading boot grammar:", err)
				os.Exit(1)
			}
			switch {
			case CLI.Boot.Json:
				output = peg.ModelToJSONStr(gram)
				lang = "json"
			case CLI.Boot.Pretty:
				output = gram.PrettyPrint()
				lang = "ebnf"
			case CLI.Boot.Railroads:
				output = gram.Railroads()
				lang = "apl"
			case CLI.Grammar.Model:
				output = util.Repr(gram)
				lang = "go"
			default:
				output = gram.PrettyPrint()
				lang = "ebnf"
			}

		case "grammar":
			gram, err := loadGrammar(CLI.Grammar.Grammar, cliCfg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			switch {
			case CLI.Grammar.Json:
				output = peg.ModelToJSONStr(gram)
				lang = "json"
			case CLI.Grammar.Pretty:
				output = gram.PrettyPrint()
				lang = "ebnf"
			case CLI.Grammar.Railroads:
				output = gram.Railroads()
				lang = "apl"
			case CLI.Grammar.Model:
				output = util.Repr(gram)
				lang = "go"
			default:
				output = gram.PrettyPrint()
				lang = "ebnf"
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
			fmt.Println(util.Pygmentize(output, lang, useColorOutput))
		}
	}
}

// loadGrammar loads a grammar from the given path, handling both EBNF and JSON formats.
func loadGrammar(path string, cfg *config.Cfg) (*peg.Grammar, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		g, err := peg.ParseGrammar(data)
		if err != nil {
			return nil, err
		}
		if err := g.Initialize(); err != nil {
			return nil, err
		}
		return g, nil
	default:
		return api.Compile(string(data), cfg)
	}
}
