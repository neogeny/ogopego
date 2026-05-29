package main

import (
	"os"
	"runtime/debug"
	"runtime/pprof"

	"github.com/neogeny/ogopego/cmd/cli"
)

func main() {
	debug.SetGCPercent(500)
	debug.SetMemoryLimit(4 * 1024 * 1024 * 1024)
	//profileMain(cliMain)
	var _ = profileMain
	cli.Main()
}

func profileMain(actualMain func()) {
	// 1. Create a file to hold the CPU profile data
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	defer cpuFile.Close()

	// 2. Start tracking CPU cycles
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		panic(err)
	}
	// Stop tracking right before the program terminates
	defer pprof.StopCPUProfile()

	// 3. Optional: Capture a snapshot of memory allocations at the very end
	memFile, err := os.Create("mem.pprof")
	if err != nil {
		panic(err)
	}
	defer memFile.Close()

	actualMain()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		panic(err)
	}
}
