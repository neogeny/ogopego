package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/neogeny/ogopego/config"
)

//goland:noinspection GoVetStructTag
var CLI struct {
	Run struct {
		Grammar string   `arg:"" required name:"grammar" help:"Path to the compiled TatSu JSON grammar."`
		Inputs  []string `arg:"" required name:"inputs" help:"The files to be parsed."`
		Json    bool     `help:"Print the output tree in JSON format" short:"j" group:"format"`
		Model   bool     `help:"Print the Rust code for the tree construction" short:"m" group:"format"`
		Short   bool     `help:"Print the Tree in short notation" short:"s" group:"format"`
	} `cmd:"" help:"Execute a grammar against one or more input files."`

	Boot struct {
		Json      bool `help:"Print the boot grammar in JSON format" short:"j" group:"format"`
		Model     bool `help:"Print the Rust code for the boot model construction" short:"m" group:"format"`
		Pretty    bool `help:"Pretty-print the boot grammar" short:"p" group:"format"`
		Railroads bool `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"boot" help:"The internal boot grammar"`

	Grammar struct {
		Grammar   string `arg:"" required name:"grammar" help:"Path to the compiled grammar (.ebnf or .json)."`
		Json      bool   `help:"Print the grammar in JSON format" short:"j" group:"format"`
		Model     bool   `help:"Print the Rust code grammar model constructors" short:"m" group:"format"`
		Pretty    bool   `help:"Pretty-print the grammar (EBNF)" short:"p" group:"format"`
		Railroads bool   `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"grammar" help:"Grammar transformations"`

	Output string `help:"Output to a file instead of stdout" short:"o"`
	Color  string `help:"Control colorized output for API results" short:"C" enum:"auto,always,never" default:"auto"`
	Trace  bool   `help:"Display a detailed trace of the parsing process." short:"t"`
}

var (
	useColorOutput bool
	cliCfg         *config.Cfg
)

func init() {
	color.NoColor = false
}

func coloredHelp(_ kong.HelpOptions, ctx *kong.Context) error {
	out := color.Output

	headerBold := color.New(color.FgYellow, color.Bold)
	flagCyan := color.New(color.FgCyan)
	cmdGreen := color.New(color.FgGreen)
	descWhite := color.New(color.FgWhite)

	selected := ctx.Selected()

	if selected == nil || ctx.Command() == "" {
		_, _ = headerBold.Fprintf(out, "Usage: ogo <command> [flags]\n\n")
		_, _ = descWhite.Fprintf(out, "ogopego: A PEG parser generator in Go\n\n")

		_, _ = flagCyan.Fprintf(out, "Flags:\n")
		_, _ = flagCyan.Fprintf(out, "  -h, --help             Show context-sensitive help.\n")
		_, _ = flagCyan.Fprintf(out, "  -o, --output=STRING    Output to a file instead of stdout\n")
		_, _ = flagCyan.Fprintf(out, "  -C, --color=auto       Control colorized output for API results\n")
		_, _ = flagCyan.Fprintf(out, "  -t, --trace           Display a detailed trace of the parsing process.\n\n")

		_, _ = cmdGreen.Fprintf(out, "Commands:\n")
		cmdGreen.Fprintf(out, "  run        Execute a grammar against one or more input files.\n")
		_, _ = cmdGreen.Fprintf(out, "  boot       The internal boot grammar\n")
		_, _ = cmdGreen.Fprintf(out, "  grammar    Grammar transformations\n\n")

		_, _ = cmdGreen.Fprintf(out, "Run \"ogo <command> --help\" for more information on a command.\n")
	} else {
		headerBold.Fprintf(out, "Usage: ogo %s %s\n\n", selected.Name, selected.FlagSummary(true))

		if selected.Detail != "" {
			_, _ = descWhite.Fprintf(out, "%s\n\n", selected.Detail)
		} else if selected.Help != "" {
			_, _ = descWhite.Fprintf(out, "%s\n\n", selected.Help)
		}

		if len(selected.Positional) > 0 {
			cmdGreen.Fprintf(out, "Arguments:\n")
			for _, arg := range selected.Positional {
				cmdGreen.Fprintf(out, "  %s\n", arg.Summary())
			}
			_, _ = out.Write([]byte("\n"))
		}

		flagCyan.Fprintf(out, "Flags:\n")
		flagCyan.Fprintf(out, "  -h, --help             Show context-sensitive help.\n")
		flagCyan.Fprintf(out, "  -o, --output=STRING    Output to a file instead of stdout\n")
		flagCyan.Fprintf(out, "  -C, --color=auto       Control colorized output for API results\n")
		flagCyan.Fprintf(out, "  -t, --trace           Display a detailed trace of the parsing process.\n")

		for _, flag := range selected.Flags {
			flagCyan.Fprintf(out, "  %s\n", flag.String())
		}
	}

	return nil
}

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
				if CLI.Boot.Model {
					set = append(set, "--model")
				}
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
