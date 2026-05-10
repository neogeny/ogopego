package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
)

var CLI struct {
	Run struct {
		Grammar string   `arg:"" required name:"grammar" help:"Path to the compiled TatSu JSON grammar."`
		Inputs  []string `arg:"" required name:"inputs" help:"The files to be parsed."`
		Json    bool    `help:"Print the output tree in JSON format" short:"j"`
		Model  bool    `help:"Print the Rust code for the tree construction" short:"m"`
		Short  bool    `help:"Print the Tree in short notation" short:"s"`
	} `cmd:"" help:"Execute a grammar against one or more input files."`

	Boot struct {
		Json      bool `help:"Print the boot grammar in JSON format" short:"j"`
		Model     bool `help:"Print the Rust code for the boot model construction" short:"m"`
		Pretty    bool `help:"Pretty-print the boot grammar" short:"p"`
		Railroads bool `help:"Print a railroad diagram" short:"r" group:"format"`
	} `cmd:"boot" help:"The internal boot grammar"`

	Grammar struct {
		Grammar   string `arg:"" required name:"grammar" help:"Path to the compiled grammar (.ebnf or .json)."`
		Json      bool   `help:"Print the grammar in JSON format" short:"j"`
		Model     bool   `help:"Print the Rust code grammar model constructors" short:"m"`
		Pretty    bool   `help:"Pretty-print the grammar (EBNF)" short:"p"`
		Railroads bool   `help:"Print a railroad diagram" short:"r"`
	} `cmd:"grammar" help:"Grammar transformations"`

	Output string `help:"Output to a file instead of stdout" short:"o"`
	Color  string `help:"Control colorized output for API results" short:"C" enum:"auto,always,never" default:"auto"`
	Trace  bool   `help:"Display a detailed trace of the parsing process." short:"t"`
}

var useColorOutput bool

func init() {
	color.NoColor = false
}

func coloredHelp(options kong.HelpOptions, ctx *kong.Context) error {
	out := color.Output

	headerBold := color.New(color.FgYellow, color.Bold)
	flagCyan := color.New(color.FgCyan)
	cmdGreen := color.New(color.FgGreen)
	descWhite := color.New(color.FgWhite)

	selected := ctx.Selected()
	
	if selected == nil || ctx.Command() == "" {
		headerBold.Fprintf(out, "Usage: ogo <command> [flags]\n\n")
		descWhite.Fprintf(out, "ogopego: A PEG parser generator in Go\n\n")
		
		flagCyan.Fprintf(out, "Flags:\n")
		flagCyan.Fprintf(out, "  -h, --help             Show context-sensitive help.\n")
		flagCyan.Fprintf(out, "  -o, --output=STRING    Output to a file instead of stdout\n")
		flagCyan.Fprintf(out, "  -C, --color=auto       Control colorized output for API results\n")
		flagCyan.Fprintf(out, "  -t, --trace           Display a detailed trace of the parsing process.\n\n")
		
		cmdGreen.Fprintf(out, "Commands:\n")
		cmdGreen.Fprintf(out, "  run        Execute a grammar against one or more input files.\n")
		cmdGreen.Fprintf(out, "  boot       The internal boot grammar\n")
		cmdGreen.Fprintf(out, "  grammar    Grammar transformations\n\n")
		
		cmdGreen.Fprintf(out, "Run \"ogo <command> --help\" for more information on a command.\n")
	} else {
		headerBold.Fprintf(out, "Usage: ogo %s %s\n\n", selected.Name, selected.FlagSummary(true))
		
		if selected.Detail != "" {
			descWhite.Fprintf(out, "%s\n\n", selected.Detail)
		} else if selected.Help != "" {
			descWhite.Fprintf(out, "%s\n\n", selected.Help)
		}

		if len(selected.Positional) > 0 {
			cmdGreen.Fprintf(out, "Arguments:\n")
			for _, arg := range selected.Positional {
				cmdGreen.Fprintf(out, "  %s\n", arg.Summary())
			}
			out.Write([]byte("\n"))
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

func isTerminal() bool {
	return os.Getenv("TERM") != "dumb"
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
		kong.Help(coloredHelp),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: false,
		}),
	)

	_ = ctx

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
}