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
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/context"
	"github.com/neogeny/ogopego/pkg/trees"
	"github.com/neogeny/ogopego/pkg/util"
)

func runCmd(cli CLIConfig, cliCfg *config.Cfg) (string, []outputItem) {
	var outputs []outputItem
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
		prog := NewProgressUI(len(cli.Run.Inputs), cli.Quiet)
		loader := prog.Loading("loading grammar")
		loadCfg := *cliCfg
		loadCfg.Heartbeat = loader.Heartbeat()
		gram, err := api.LoadGrammar(cli.Run.Grammar, &loadCfg)
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
				prog.IncFiles()
				<-sem

				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					fp.Fail()
					errcount++
					if report, ok := errors.AsType[*context.ParseFailure](err); ok {
						err = &report.Memento
					}
					_, _ = fmt.Fprintf(os.Stderr, "\n%v\n", err)
					return
				}

				sourceLines += util.CountLines(input).Code
				fp.Success()
				var payload string
				switch {
				case cli.Run.Model:
					payload = util.Repr(tree)
				default:
					payload = trees.TreeToJSONStr(tree)
				}
				item := outputItem{Name: fName, Payload: payload}
				outputs = append(outputs, item)
			}(path)
		}

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

	var lang string
	switch {
	case cli.Run.Model:
		lang = "go"
	default:
		lang = "json"
	}
	return lang, outputs
}
