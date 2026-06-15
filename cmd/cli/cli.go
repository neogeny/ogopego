// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package cli

// CLIConfig holds the parsed command-line configuration for the ogo tool.
type CLIConfig struct {
	Run struct {
		Grammar string   `arg:"required,name=grammar" help:"Path to the grammar in EBNF or JSON format"`
		Inputs  []string `arg:"required,name=inputs" help:"The files to be parsed"`
		Json    bool     `help:"Print the output tree in JSON format" short:"j" xor:"format"`
		Jsonl   bool     `help:"Output JSON Lines (compact, one object per line)"`
		Model   bool     `help:"Print the Go code for the tree construction" short:"m" xor:"format"`
		Start   string   `help:"Name of the start rule (defaults to 'start')" short:"s"`
		Nproc   int      `help:"Number of concurrent workers (default: number of CPUs)" short:"n" default:"0"`
	} `cmd:"" help:"Execute a grammar against one or more input files"`

	Boot struct {
		Json      bool `help:"Print the boot grammar in JSON format" short:"j" xor:"format"`
		Model     bool `help:"Print the Go code for the boot model construction" short:"m" xor:"format"`
		Pretty    bool `help:"Pretty-print the boot grammar" short:"p" xor:"format"`
		Railroads bool `help:"Print a railroad diagram" short:"r" xor:"format"`
	} `cmd:"boot" help:"The internal boot grammar"`

	Grammar struct {
		Grammar   string `arg:"required,name=grammar" help:"Path to the grammar source (.ebnf or .json)"`
		Json      bool   `help:"Print the grammar in JSON format" short:"j" xor:"format"`
		Model     bool   `help:"Print the Go code grammar model constructors" short:"m" xor:"format"`
		Parser    string `help:"Generate Go parser source code" short:"x" xor:"format" placeholder:"PKG"`
		ModelGen  string `help:"Generate Go model source code" short:"g" xor:"format" placeholder:"PKG"`
		Pretty    bool   `help:"Pretty-print the grammar (EBNF)" short:"p" xor:"format"`
		Railroads bool   `help:"Print a railroad diagram" short:"r" xor:"format"`
	} `cmd:"grammar" help:"Grammar transformations"`

	Profile bool   `help:"Enable CPU and memory profiling, output to $TMPDIR"`
	Output  string `help:"Output to a file or directory instead of stdout" short:"o"`
	Color   string `help:"Control colorized output for API results" enum:"auto,always,never" default:"auto"`
	Trace   bool   `help:"Display a detailed trace of the parsing process"`
	Verbose bool   `help:"Display full error context with source lines" short:"v"`
	Quiet   bool   `help:"Suppress progress bar and spinner output"`
	Version bool   `help:"Print version information"`
}
