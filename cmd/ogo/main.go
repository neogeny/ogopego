package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
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
		fmt.Fprintln(os.Stderr, err)
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

	switch cmd.Name {
	case "run":
		fmt.Fprintln(os.Stderr, "grammar:", CLI.Run.Grammar)
		fmt.Fprintln(os.Stderr, "inputs:", CLI.Run.Inputs)
		fmt.Fprintln(os.Stderr, "json:", CLI.Run.Json)
		fmt.Fprintln(os.Stderr, "model:", CLI.Run.Model)
		fmt.Fprintln(os.Stderr, "short:", CLI.Run.Short)
	case "boot":
		fmt.Fprintln(os.Stderr, "json:", CLI.Boot.Json)
		fmt.Fprintln(os.Stderr, "model:", CLI.Boot.Model)
		fmt.Fprintln(os.Stderr, "pretty:", CLI.Boot.Pretty)
		fmt.Fprintln(os.Stderr, "railroads:", CLI.Boot.Railroads)
	case "grammar":
		fmt.Fprintln(os.Stderr, "grammar:", CLI.Grammar.Grammar)
		fmt.Fprintln(os.Stderr, "json:", CLI.Grammar.Json)
		fmt.Fprintln(os.Stderr, "model:", CLI.Grammar.Model)
		fmt.Fprintln(os.Stderr, "pretty:", CLI.Grammar.Pretty)
		fmt.Fprintln(os.Stderr, "railroads:", CLI.Grammar.Railroads)
	}
}
