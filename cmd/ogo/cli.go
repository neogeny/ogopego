// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/neogeny/ogopego/config"
)

//goland:noinspection GoVetStructTag
var CLI struct {
	Run struct {
		Grammar string   `arg:"" required name:"grammar" help:"Path to the grammar in EBNF or JSON format"`
		Inputs  []string `arg:"" required name:"inputs" help:"The files to be parsed"`
		Json    bool     `help:"Print the output tree in JSON format" short:"j" group:"format"`
		//Model   bool     `help:"Print the Go code for the tree construction" short:"m" group:"format"`
		Short bool `help:"Print the Tree in short notation" short:"s" group:"format"`
	} `cmd:"" help:"Execute a grammar against one or more input files"`

	Boot struct {
		Json bool `help:"Print the boot grammar in JSON format" short:"j" group:"format"`
		//Model     bool `help:"Print the Go code for the boot model construction" short:"m" group:"format"`
		Pretty    bool `help:"Pretty-print the boot grammar" short:"p" group:"format"`
		Railroads bool `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"boot" help:"The internal boot grammar"`

	Grammar struct {
		Grammar string `arg:"" required name:"grammar" help:"Path to the compiled grammar (.ebnf or .json)"`
		Json    bool   `help:"Print the grammar in JSON format" short:"j" group:"format"`
		//Model     bool   `help:"Print the Go code grammar model constructors" short:"m" group:"format"`
		Pretty    bool `help:"Pretty-print the grammar (EBNF)" short:"p" group:"format"`
		Railroads bool `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"grammar" help:"Grammar transformations"`

	Output string `help:"Output to a file instead of stdout" short:"o"`
	Color  string `help:"Control colorized output for API results" short:"C" enum:"auto,always,never" default:auto`
	Trace  bool   `help:"Display a detailed trace of the parsing process" short:"t"`
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
				//if CLI.Run.Model {
				//	set = append(set, "--model")
				//}
				if CLI.Run.Short {
					set = append(set, "--short")
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
				//if CLI.Grammar.Model {
				//	set = append(set, "--model")
				//}
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
		return fmt.Errorf("only one of --json, --model, --pretty, --railroads, --short can be specified")
	}
	return nil
}

func isTerminal() bool {
	return os.Getenv("TERM") != "dumb"
}

var ctx *kong.Context
