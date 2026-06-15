package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/util"
)

var (
	summaryFilesStyle  = color.New(color.FgWhite, color.Bold)
	summaryPassedStyle = color.New(color.FgGreen)
	summaryFailStyle   = color.New(color.FgRed)
	summaryRateStyle   = color.New(color.FgCyan)

	tableLabelStyle = color.New(color.FgCyan, color.Faint)
	tableValueStyle = color.New(color.FgWhite, color.Bold)
	tableGoodStyle  = color.New(color.FgGreen)
	tableBadStyle   = color.New(color.FgRed)
	tableMidStyle   = color.New(color.FgYellow)

	diagErrStyle = color.New(color.FgRed)
)

func runCmd(cli CLIConfig, cliCfg *config.Cfg) (string, []outputItem) {
	var outputs []outputItem
	var errcount int
	var totlLines int
	var codeLines int
	var cmntLines int
	var blnkLines int
	var succCount int
	var succLines int
	var runTime float64
	startTime := time.Now()

	if cli.Run.Start != "" {
		cliCfg.Start = cli.Run.Start
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
		loadCfg := *cliCfg
		loadCfg.Heart = loader.Heartbeat()
		gram, err := api.LoadGrammar(cli.Run.Grammar, &loadCfg)
		loader.Finish()
		if err != nil {
			_, _ = diagErrStyle.Fprintln(os.Stderr, "\nerror:", err)
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
				fileStart := time.Now()
				lc := util.CountLines(input)
				fp.SetLength(len(data))

				fileCfg := *cliCfg
				fileCfg.Heart = fp.Heartbeat()
				fileCfg.Source, _ = util.PathRelativeToCwd(path)

				tree, err := api.ParseInput(gram, input, &fileCfg)
				prog.IncFiles()
				<-sem

				mu.Lock()
				defer mu.Unlock()
				totlLines += lc.Total
				codeLines += lc.Code
				cmntLines += lc.Comment
				blnkLines += lc.Blank
				runTime += time.Since(fileStart).Seconds()
				if err != nil {
					fp.Fail()
					errcount++
					if report, ok := errors.AsType[*context.ParseFailure](err); ok {
						err = &report.Memento
					}
					_, _ = fmt.Fprintf(os.Stderr, "\n%v\n", err)
					return
				}

				succCount++
				succLines += lc.Total
				fp.Success()
				var payload string
				switch {
				case cli.Run.Model:
					payload = util.Repr(tree)
				default:
					payload = asjson.AsJSONStr(tree)
				}
				item := outputItem{Name: fName, Payload: payload}
				outputs = append(outputs, item)
			}(path)
		}

		wg.Wait()
		prog.Finish()
	}

	if runTime > 0 {
		slocsAvg := float64(totlLines) / runTime
		succRate := float64(succCount) / float64(len(inputs))
		failures := len(inputs) - succCount
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
			{"       files input", fmt.Sprintf("%d", len(inputs)), tableLabelStyle, tableValueStyle},
			{" source lines input", fmt.Sprintf("%d", totlLines), tableLabelStyle, tableValueStyle},
			{"     success lines", fmt.Sprintf("%d", succLines), tableLabelStyle, tableValueStyle},
			{"              sloc", fmt.Sprintf("%d", codeLines), tableLabelStyle, tableValueStyle},
			{"         successes", fmt.Sprintf("%d", succCount), tableGoodStyle, tableGoodStyle},
			{"          failures", fmt.Sprintf("%d", failures), tableBadStyle, tableBadStyle},
			{"      success rate", fmt.Sprintf("%12.0f %%", 100.0*succRate), tableLabelStyle, rateClr},
			{"         sloc/sec", fmt.Sprintf("%12.0f sl/s", slocsAvg), tableLabelStyle, slocClr},
			{"          run time", fmtDuration(runTime), tableLabelStyle, tableValueStyle},
			{"         wall time", fmtDuration(wallTime), tableLabelStyle, tableValueStyle},
		}

		for _, r := range rows {
			r.labelClr.Fprintf(os.Stderr, "%20s ", r.label)
			r.valueClr.Fprintf(os.Stderr, "%12s\n", r.value)
		}
	}

	/*
		fmt.Fprintln(os.Stderr)

		passed := len(cli.Run.Inputs) - errcount
		elapsed := time.Since(startTime)
		rate := int(float64(codeLines) / elapsed.Seconds())

		errors := ""
		if errcount > 0 {
			errors = summaryFailStyle.Sprintf(" %d errors", errcount)
		}
		_, _ = fmt.Fprintf(os.Stderr, "Parsed%s%s%s%s\n",
			summaryFilesStyle.Sprintf(" %d files", len(cli.Run.Inputs)),
			summaryPassedStyle.Sprintf(" %d passed", passed),
			errors,
			summaryRateStyle.Sprintf(" %d sloc/s", rate),
		)
	*/

	var lang string
	switch {
	case cli.Run.Model:
		lang = "go"
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
