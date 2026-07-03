// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/parproc"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/util"
)

var (
	tableLabelStyle = color.New(color.FgCyan, color.Faint)
	tableValueStyle = color.New(color.FgWhite, color.Bold)
	tableGoodStyle  = color.New(color.FgGreen)
	tableBadStyle   = color.New(color.FgRed)
	tableMidStyle   = color.New(color.FgYellow)
	tableDimStyle   = color.New(color.FgWhite, color.Faint)

	diagErrStyle = color.New(color.FgRed)
)

type tPayload struct {
	Path string
	Text string
	Gram *peg.Grammar
	Cfg  *config.Cfg
	Prog *ProgressUI
}

func runCmd(cli CLIConfig, parserConfig *config.Cfg) (string, []tOutputItem) {
	var outputs []tOutputItem
	var parseFailures []error
	var totlLines int
	var codeLines int
	var cmntLines int
	var blnkLines int
	var succCount int
	var failCount int
	var succLines int
	var runTime float64
	startTime := time.Now()

	if cli.Run.Start != "" {
		parserConfig.Start = cli.Run.Start
	}

	var inputs []string
	for _, path := range cli.Run.Inputs {
		if util.FileExists(path) {
			inputs = append(inputs, path)
		} else {
			_, _ = diagErrStyle.
				Fprintf(os.Stderr, "warning: input file not found: %s\n", path)
		}
	}

	if len(inputs) > 0 {
		prog := NewProgressUI(len(cli.Run.Inputs), cli.Quiet)
		loader := prog.Loading("loading grammar")
		loadCfg := *parserConfig
		loadCfg.Heart = loader.Heartbeat()
		gram, err := api.LoadGrammar(cli.Run.Grammar, &loadCfg)
		loader.Finish()
		if err != nil {
			_, _ = diagErrStyle.Fprintln(os.Stderr, "\nerror:", err)
			os.Exit(1)
		}

		maxWorkers := cli.Run.Nproc
		var payloads []tPayload
		for _, path := range inputs {
			text, err := os.ReadFile(path)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\nerror reading %s: %v\n", path, err)
				continue
			}
			payloads = append(
				payloads,
				tPayload{
					Path: path,
					Text: string(text),
					Gram: gram,
					Cfg:  parserConfig,
					Prog: prog,
				},
			)

		}
		for result := range parproc.ParProc(parseTask, payloads, maxWorkers) {
			payload := result.Payload
			name := filepath.Base(payload.Path)
			tree := result.Outcome

			lc := util.CountLines(payload.Text)
			totlLines += lc.Total
			codeLines += lc.Code
			cmntLines += lc.Comment
			blnkLines += lc.Blank
			runTime += result.Elapsed.Seconds()

			if cli.Verbose && !cli.Quiet {
				var style *color.Color
				var symbol string
				if result.Error != nil {
					style = tableBadStyle
					symbol = "✗"
				} else {
					style = tableGoodStyle
					symbol = "✓"
				}
				msg := style.Sprintf("%3s", symbol) +
					tableDimStyle.Sprintf(" %-50s ", name) +
					style.Sprintf("⏲ %6.2f s\n", result.Elapsed.Seconds())
				_, _ = prog.Write([]byte(msg))
			}
			if result.Error != nil {
				err := result.Error
				if report, ok := errors.AsType[*context.ParseFailure](err); ok {
					err = &report.Memento
				}
				parseFailures = append(parseFailures, result.Error)
				continue
			}

			succCount++
			succLines += lc.Total
			var output string
			switch {
			case cli.Run.Model:
				output = util.Repr(tree)
			default:
				output = asjson.AsJSONStr(tree)
			}
			outputs = append(outputs, tOutputItem{
				Path:   payload.Path,
				Output: output,
			})
		}
		prog.Finish()
	}

	if !cli.Quiet && runTime > 0 {
		total := len(inputs)
		failCount = total - succCount
		slocsAvg := float64(totlLines) / runTime
		succRate := float64(succCount) / float64(total)
		wallTime := time.Since(startTime).Seconds()

		fmt.Fprintln(os.Stderr)
		type tableRow struct {
			label    string
			value    string
			labelClr *color.Color
			valueClr *color.Color
		}

		var rateClr *color.Color
		switch {
		case succRate >= 1.0:
			rateClr = tableGoodStyle
		case succRate > 0.6:
			rateClr = tableMidStyle
		default:
			rateClr = tableBadStyle
		}

		var slocClr *color.Color
		switch {
		case slocsAvg >= 200:
			slocClr = tableGoodStyle
		case slocsAvg >= 180:
			slocClr = tableMidStyle
		default:
			slocClr = tableBadStyle
		}

		rows := []tableRow{
			{"       files input", fmt.Sprintf("%d", total), tableLabelStyle, tableValueStyle},
			{" source lines input", fmt.Sprintf("%d", totlLines), tableLabelStyle, tableValueStyle},
			{"     success lines", fmt.Sprintf("%d", succLines), tableLabelStyle, tableValueStyle},
			{"              sloc", fmt.Sprintf("%d", codeLines), tableLabelStyle, tableValueStyle},
			{"         successes", fmt.Sprintf("%d", succCount), tableGoodStyle, tableGoodStyle},
			{"          failures", fmt.Sprintf("%d", failCount), tableBadStyle, tableBadStyle},
			{"      success rate", fmt.Sprintf("%12.0f %%", 100.0*succRate), tableLabelStyle, rateClr},
			{"         sloc/sec", fmt.Sprintf("%12.0f sl/s", slocsAvg), tableLabelStyle, slocClr},
			{"          run time", fmtDuration(runTime), tableLabelStyle, tableValueStyle},
			{"         wall time", fmtDuration(wallTime), tableLabelStyle, tableValueStyle},
		}

		if cli.Verbose {
			for _, err := range parseFailures {
				if report, ok := errors.AsType[*context.ParseFailure](err); ok {
					err = &report.Memento
				}
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}

		for _, r := range rows {
			r.labelClr.Fprintf(os.Stderr, "%20s ", r.label)
			r.valueClr.Fprintf(os.Stderr, "%12s\n", r.value)
		}
	}

	var lang string
	switch {
	case cli.Run.Model:
		lang = "go"
	case cli.Run.Jsonl:
		lang = "jsonl"
	default:
		lang = "json"
	}
	return lang, outputs
}

func fmtDuration(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second)).Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h >= 1 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}
func parseTask(payload tPayload) (any, error) {
	path := payload.Path
	fName := filepath.Base(path)
	fp := payload.Prog.AddFile(fName)
	fp.SetLength(len(payload.Text))

	fileCfg := payload.Cfg
	fileCfg.Heart = fp.Heartbeat()
	fileCfg.Source, _ = util.PathRelativeToCwd(path)

	tree, err := api.ParseInput(payload.Gram, payload.Text, fileCfg)
	payload.Prog.IncFiles()
	if err != nil {
		fp.Fail()
		return nil, err
	}
	fp.Success()
	return tree, nil
}
