// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/util"
)

const (
	modeStdout = iota
	modeFile
	modeDir
)

type outputItem struct {
	Name    string
	Payload string
}

func parseCLI() (CLIConfig, *kong.Context) {
	var cfg CLIConfig
	ctx := kong.Parse(&cfg,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: false,
		}),
	)
	return cfg, ctx
}

func Main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			fmt.Printf("%s %s\n", config.ProgramName, util.GetVersion())
			os.Exit(0)
		}
	}

	cli, ctx := parseCLI()

	var (
		useColorOutput bool
		outputs        []outputItem
		lang           string
	)

	switch cli.Color {
	case "always":
		useColorOutput = true
		color.NoColor = false
	case "never":
		useColorOutput = false
		color.NoColor = true
	case "auto":
		if util.IsTerminal() {
			useColorOutput = true
			color.NoColor = false
		} else {
			useColorOutput = false
			color.NoColor = true
		}
	}

	cliCfg := &config.Cfg{
		Trace:    cli.Trace,
		Colorize: useColorOutput,
	}

	cmd := ctx.Selected()
	if cmd == nil {
		return
	}

	switch cmd.Name {
	case "run":
		lang, outputs = runCmd(cli, cliCfg)

	case "boot":
		lang, outputs = bootCmd(cli)

	case "grammar":
		lang, outputs = grammarCmd(cli, cliCfg)
	}

	if len(outputs) > 0 {
		if err := writeOutputs(outputs, lang, cli.Output, useColorOutput); err != nil {
			fmt.Fprintln(os.Stderr, "error writing output:", err)
			os.Exit(1)
		}
	}
}

func outputMode(path string) int {
	if path == "" || path == "-" {
		return modeStdout
	}
	if path == "/dev/null" {
		return modeFile
	}
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			return modeDir
		}
		return modeFile
	}
	if filepath.Ext(path) == "" {
		return modeDir
	}
	return modeFile
}

func langExt(lang string) string {
	switch lang {
	case "json":
		return ".json"
	case "jsonl":
		return ".jsonl"
	case "go":
		return ".go"
	default:
		return ".ebnf"
	}
}

func replaceExt(name, newExt string) string {
	if old := filepath.Ext(name); old != "" {
		name = name[:len(name)-len(old)]
	}
	return name + newExt
}

func writeOutputs(outputs []outputItem, lang string, path string, color bool) error {
	switch outputMode(path) {
	case modeStdout:
		if lang == "jsonl" {
			out, _ := formatOutputs(outputs, true)
			_, _ = os.Stdout.Write(out)
		} else {
			fmt.Println(Pygmentize(joinOutputs(outputs), lang, color))
		}
		return nil

	case modeFile:
		out, err := formatOutputs(outputs, lang == "jsonl")
		if err != nil {
			return err
		}
		return os.WriteFile(path, out, 0644)

	case modeDir:
		ext := langExt(lang)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
		for _, o := range outputs {
			outPath := filepath.Join(path, replaceExt(o.Name, ext))
			data := []byte(o.Payload)
			if lang == "jsonl" {
				var buf bytes.Buffer
				if err := json.Compact(&buf, data); err != nil {
					return fmt.Errorf("compacting json for %s: %w", o.Name, err)
				}
				buf.WriteByte('\n')
				data = buf.Bytes()
			}
			if err := os.WriteFile(outPath, data, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}
		}
		return nil
	}
	return nil
}

func formatOutputs(outputs []outputItem, compact bool) ([]byte, error) {
	if compact {
		var buf bytes.Buffer
		for _, o := range outputs {
			if err := json.Compact(&buf, []byte(o.Payload)); err != nil {
				return nil, fmt.Errorf("compacting json: %w", err)
			}
			buf.WriteByte('\n')
		}
		return buf.Bytes(), nil
	}
	return []byte(joinOutputs(outputs)), nil
}

func joinOutputs(outputs []outputItem) string {
	payloads := make([]string, len(outputs))
	for i, o := range outputs {
		payloads[i] = o.Payload
	}
	return strings.Join(payloads, "\n")
}
