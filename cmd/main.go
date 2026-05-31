package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/neogeny/ogopego/cmd/cli"
)

func main() {
	debug.SetGCPercent(500)
	debug.SetMemoryLimit(4 * 1024 * 1024 * 1024)

	if profileEnabled() {
		profileMain(cli.Main)
	} else {
		cli.Main()
	}
}

func profileEnabled() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--profile" || strings.HasPrefix(arg, "--profile=") {
			return true
		}
	}
	return false
}

func profileMain(actualMain func()) {
	timestamp := time.Now().Format("2006-01-02-1504")

	cpuPath := fmt.Sprintf("cpu-%s.pprof", timestamp)
	cpuFile, err := os.Create(cpuPath)
	if err != nil {
		panic(err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	runtime.MemProfileRate = 16 * 1024

	actualMain()

	memPath := fmt.Sprintf("mem-%s.pprof", timestamp)
	memFile, err := os.Create(memPath)
	if err != nil {
		panic(err)
	}
	defer memFile.Close()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		panic(err)
	}

	allocPath := fmt.Sprintf("alc-%s.pprof", timestamp)
	allocFile, err := os.Create(allocPath)
	if err != nil {
		panic(err)
	}
	defer allocFile.Close()

	if err := pprof.Lookup("allocs").WriteTo(allocFile, 0); err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "\nprofiles:\n  %s\n  %s\n  %s\n", cpuPath, memPath, allocPath)
}
