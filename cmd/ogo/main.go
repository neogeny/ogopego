package main

import (
	"os"

	"github.com/alecthomas/kong"
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
	Color  string `help:"Control when to use color in output." short:"C" enum:"auto,always,never" default:"auto"`
	Trace  bool   `help:"Display a detailed trace of the parsing process." short:"t"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
	)

	_ = ctx
	_ = os.Stdout
}