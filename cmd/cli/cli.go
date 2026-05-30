// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/neogeny/ogopego/pkg/config"
)

var ctx *kong.Context

func IsTerminal() bool {
	return os.Getenv("TERM") != "dumb"
}

// CLI is the command-line interface structure for the ogo tool.
//
//goland:noinspection GoVetStructTag
//nolint:govet // Kong uses non-standard tag syntax
var CLI struct {
	// Run subcommand executes a grammar against one or more input files.
	Run struct {
		Grammar string   `arg required name:"grammar" help:"Path to the grammar in EBNF or JSON format"`
		Inputs  []string `arg required name:"inputs" help:"The files to be parsed"`
		Json    bool     `help:"Print the output tree in JSON format" short:"j" group:"format"`
		Model   bool     `help:"Print the Go code for the tree construction" short:"m" group:"format"`
		Start   string   `help:"Name of the start rule (defaults to 'start')" short:"s"`
	} `cmd:"" help:"Execute a grammar against one or more input files"`

	// Boot subcommand provides access to the internal boot grammar.
	Boot struct {
		Json      bool `help:"Print the boot grammar in JSON format" short:"j" group:"format"`
		Model     bool `help:"Print the Go code for the boot model construction" short:"m" group:"format"`
		Pretty    bool `help:"Pretty-print the boot grammar" short:"p" group:"format"`
		Railroads bool `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"boot" help:"The internal boot grammar"`

	// Grammar subcommand provides transformations and information about a grammar.
	Grammar struct {
		Grammar   string `arg required name:"grammar" help:"Path to the grammar source (.ebnf or .json)"`
		Json      bool   `help:"Print the grammar in JSON format" short:"j" group:"format"`
		Model     bool   `help:"Print the Go code grammar model constructors" short:"m" group:"format"`
		Parser    string `help:"Generate Go parser source code" short:"x" group:"format" placeholder:"PKG"`
		ModelGen  string `help:"Generate Go model source code" short:"g" group:"format" placeholder:"PKG"`
		Pretty    bool   `help:"Pretty-print the grammar (EBNF)" short:"p" group:"format"`
		Railroads bool   `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"grammar" help:"Grammar transformations"`

	// Output specifies a file to write output to instead of stdout.
	Output string `help:"Output to a file or directory instead of stdout" short:"o"`
	// Color controls colorized output for API results.
	Color string `help:"Control colorized output for API results" enum:"auto,always,never" default:"auto"`
	// Trace enables a detailed trace of the parsing process.
	Trace bool `help:"Display a detailed trace of the parsing process"`
	// Quiet disables progress bars and spinner output.
	Quiet bool `help:"Suppress progress bar and spinner output"`
	// Version prints version information.
	Version bool `help:"Print version information"`
}

var (
	useColorOutput bool
	cliCfg         *config.Cfg
)

func validateExclusive(groups ...string) error {
	cmd := ctx.Selected()
	if cmd == nil {
		return nil
	}
	selected := cmd.Name
	var set []string
	for _, group := range groups {
		switch selected {
		case "run":
			switch group {
			case "format":
				if CLI.Run.Json {
					set = append(set, "--json")
				}
				if CLI.Run.Model {
					set = append(set, "--model")
				}
			}
		case "boot":
			switch group {
			case "format":
				if CLI.Boot.Json {
					set = append(set, "--json")
				}
				//if CLI.Boot.Model {
				//	set = append(set, "--model")
				//}
				if CLI.Boot.Pretty {
					set = append(set, "--pretty")
				}
				if CLI.Boot.Railroads {
					set = append(set, "--railroads")
				}
			}
		case "grammar":
			switch group {
			case "format":
				if CLI.Grammar.Json {
					set = append(set, "--json")
				}
				if CLI.Grammar.Model {
					set = append(set, "--model")
				}
				if CLI.Grammar.Parser != "" {
					set = append(set, "--parser="+CLI.Grammar.Parser)
				}
				if CLI.Grammar.ModelGen != "" {
					set = append(set, "--model-gen="+CLI.Grammar.ModelGen)
				}
				if CLI.Grammar.Pretty {
					set = append(set, "--pretty")
				}
				if CLI.Grammar.Railroads {
					set = append(set, "--railroads")
				}
			}
		}
	}
	if len(set) > 1 {
		return fmt.Errorf("only one of --json, --model, --parser, --pretty, --railroads can be specified")
	}
	return nil
}
