package main

import (
	"os"
	"runtime/pprof"
)

func main() {
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

	cliMain()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		panic(err)
	}
}
