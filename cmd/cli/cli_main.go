// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/tool"
	"github.com/neogeny/ogopego/pkg/trees"
	"github.com/neogeny/ogopego/pkg/util"
	"github.com/vbauerster/mpb/v8"
)

type outputItem struct {
	Name    string
	Payload string
}

func Main() {
	var (
		cliCfg         *config.Cfg
		useColorOutput bool
		outputs        []outputItem
		lang           string
	)

	var ctx *kong.Context
	ctx = kong.Parse(&cli,
		kong.Name("ogo"),
		kong.Description("ogopego: A PEG parser generator in Go"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: false,
		}),
	)

	if cli.Version {
		fmt.Printf("%s %s\n", config.ProgramName, util.GetVersion())
		os.Exit(0)
	}

	if err := validateExclusive(ctx, "format"); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	cliCfg = &config.Cfg{
		Trace:    cli.Trace,
		Colorize: useColorOutput,
	}

	cmd := ctx.Selected()
	if cmd == nil {
		return
	}
	switch cmd.Name {
	case "run":
		var errcount int
		var sourceLines int
		start := time.Now()

		if cli.Run.Start != "" {
			cliCfg.Start = cli.Run.Start
		}
		var inputs []string
		for _, path := range cli.Run.Inputs {
			if util.FileExists(path) {
				inputs = append(inputs, path)
			} else {
				_, _ = color.New(color.FgRed).
					Fprintf(os.Stderr, "warning: input file not found: %s\n", path)
			}
		}

		if len(inputs) > 0 {
			//panic(fmt.Sprintf("HERE %v", CLI.Run.Inputs))
			prog := NewProgressUI(len(cli.Run.Inputs), cli.Quiet)
			loader := prog.Loading("loading grammar")
			loadCfg := *cliCfg
			loadCfg.Heartbeat = loader.Heartbeat()
			gram, err := loadGrammar(cli.Run.Grammar, &loadCfg)
			loader.Finish()
			if err != nil {
				_, _ = color.New(color.FgRed).Fprintln(os.Stderr, "\nerror:", err)
				os.Exit(1)
			}

			maxWorkers := cli.Run.Nproc
			if maxWorkers <= 0 {
				maxWorkers = runtime.GOMAXPROCS(0)
			}
			var mu sync.Mutex
			var wg sync.WaitGroup
			var sem = make(chan int, maxWorkers)
			for i, path := range inputs {
				wg.Add(1)
				sem <- i
				go func(path string) {
					defer wg.Done()
					fName := filepath.Base(path)
					fp := prog.AddFile(fName)

					data, err := os.ReadFile(path)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "\nerror reading %s: %v\n", path, err)
						errcount++
						prog.IncFiles()
						return
					}
					input := string(data)
					fp.SetLength(len(data))

					fileCfg := *cliCfg
					fileCfg.Heartbeat = fp.Heartbeat()
					fileCfg.Source, _ = util.PathRelativeToCwd(path)

					tree, err := api.ParseInput(gram, input, &fileCfg)
					// Thread-safe accumulation block
					prog.IncFiles()
					<-sem
					mu.Lock()
					defer mu.Unlock()

					if err != nil {
						errcount++
						fp.Fail()
						if report, ok := errors.AsType[*context.ParseFailure](err); ok {
							err = &report.Memento
						}
						_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
					} else {
						sourceLines += util.CountLines(input).Code
						fp.Success()
						var payload string
						switch {
						case cli.Run.Model:
							payload = util.Repr(tree)
						default:
							payload = trees.TreeToJSONStr(tree)
							item := outputItem{Name: fName, Payload: payload}
							outputs = append(outputs, item)
						}
					}
				}(path)
			}

			// Block the main thread until every individual parsing worker finishes
			wg.Wait()
			prog.Finish()
		}
		passed := len(cli.Run.Inputs) - errcount
		elapsed := time.Since(start)
		rate := int(float64(sourceLines) / elapsed.Seconds())
		_, _ = fmt.Fprintf(os.Stderr, "%s %s %s %s\n",
			color.New(color.FgWhite, color.Bold).Sprintf("Parsed %d files", len(cli.Run.Inputs)),
			color.New(color.FgGreen, color.Bold).Sprintf("%d passed", passed),
			color.New(color.FgRed, color.Bold).Sprintf("%d errors", errcount),
			color.New(color.FgCyan).Sprintf("%d sloc/s", rate),
		)
		switch {
		case cli.Run.Model:
			lang = "go"
		default:
			lang = "json"
		}

	case "boot":
		gram, err := api.BootGrammar()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error loading boot grammar:", err)
			os.Exit(1)
		}
		var payload string
		switch {
		case cli.Boot.Json:
			payload = peg.ModelToJSONStr(gram)
			lang = "json"
		case cli.Boot.Pretty:
			payload = gram.PrettyPrint()
			lang = "ebnf"
		case cli.Boot.Railroads:
			payload = gram.Railroads()
			lang = "apl"
		case cli.Boot.Model:
			payload = util.Repr(gram)
			lang = "go"
		default:
			payload = gram.PrettyPrint()
			lang = "ebnf"
		}
		outputs = append(outputs, outputItem{Name: "boot", Payload: payload})

	case "grammar":
		var p *mpb.Progress
		if !cli.Quiet {
			p = mpb.New(mpb.WithOutput(os.Stderr))
		}
		fileName := filepath.Base(cli.Grammar.Grammar)
		fp := NewFileProgress(p, fileName)

		data, err := os.ReadFile(cli.Grammar.Grammar)
		if err != nil {
			fmt.Fprintln(os.Stderr, "\nerror reading grammar:", err)
			os.Exit(1)
		}
		fp.SetLength(len(data))

		loadCfg := *cliCfg
		loadCfg.Heartbeat = fp.Heartbeat()
		gram, err := loadGrammar(cli.Grammar.Grammar, &loadCfg)
		if err != nil {
			fp.Fail()
			if p != nil {
				p.Wait()
			}
			fmt.Fprintln(os.Stderr, "\nerror:", err)
			os.Exit(1)
		}
		fp.Success()
		if p != nil {
			p.Wait()
		}
		var payload string
		switch {
		case cli.Grammar.Json:
			payload = peg.ModelToJSONStr(gram)
			lang = "json"
		case cli.Grammar.Pretty:
			payload = gram.PrettyPrint()
			lang = "ebnf"
		case cli.Grammar.Railroads:
			payload = gram.Railroads()
			lang = "apl"
		case cli.Grammar.Model:
			payload = util.Repr(gram)
			lang = "go"
		case cli.Grammar.Parser != "":
			payload = peg.ParserRepr(*gram, cli.Grammar.Parser)
			lang = "go"
		case cli.Grammar.ModelGen != "":
			payload = tool.ModelRepr(*gram, cli.Grammar.ModelGen)
			lang = "go"
		default:
			payload = gram.PrettyPrint()
			lang = "ebnf"
		}
		outputs = append(outputs, outputItem{Name: filepath.Base(cli.Grammar.Grammar), Payload: payload})
	}

	if len(outputs) > 0 {
		if err := writeOutputs(outputs, lang, cli.Output, useColorOutput); err != nil {
			fmt.Fprintln(os.Stderr, "error writing output:", err)
			os.Exit(1)
		}
	}
}

// outputMode classifies an output path as stdout, file, or directory mode.
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

const (
	modeStdout = iota
	modeFile
	modeDir
)

// langExt returns the file extension for a given output language.
func langExt(lang string) string {
	switch lang {
	case "json":
		return ".json"
	case "go":
		return ".go"
	default:
		return ".ebnf"
	}
}

// replaceExt replaces the extension of name with newExt.
// If name has no extension, newExt is appended.
func replaceExt(name, newExt string) string {
	if old := filepath.Ext(name); old != "" {
		name = name[:len(name)-len(old)]
	}
	return name + newExt
}

// writeOutputs routes outputs to stdout, a single file, or a directory.
func writeOutputs(outputs []outputItem, lang string, path string, color bool) error {
	switch outputMode(path) {
	case modeStdout:
		fmt.Println(Pygmentize(joinOutputs(outputs), lang, color))
		return nil

	case modeFile:
		joined := joinOutputs(outputs)
		return os.WriteFile(path, []byte(joined), 0644)

	case modeDir:
		ext := langExt(lang)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
		for _, o := range outputs {
			name := replaceExt(o.Name, ext)
			outPath := filepath.Join(path, name)
			if err := os.WriteFile(outPath, []byte(o.Payload), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}
		}
		return nil
	}
	return nil
}

func joinOutputs(outputs []outputItem) string {
	payloads := make([]string, len(outputs))
	for i, o := range outputs {
		payloads[i] = o.Payload
	}
	return strings.Join(payloads, "\n")
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
		g, err := peg.LoadGrammarFromJSON(data)
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
