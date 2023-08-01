package main

import (
	"flag"
	"fmt"
	"time"
)

func measureExecutionTime(function func()) time.Duration {
	startTime := time.Now()
	function()
	endTime := time.Now()
	return endTime.Sub(startTime)
}

func main() {
	filePath := *flag.String("file", "prt_fild02.fld2", "filepath to the .fld2 file")

	grid, err := ParseGridFromFile(filePath)
	if err != nil {
		fmt.Println("Error parsing the file:", err)
		return
	}

	// Precompute all paths and store them in a cache.
	cache := PrecomputePaths(grid)

	// Dump the paths cache into a JSON file.
	jsonFilePath := "paths_cache.json"
	err = SavePathsCacheToJSON(cache, jsonFilePath)
	if err != nil {
		fmt.Println("Error saving paths cache to JSON:", err)
		return
	}

	fmt.Println("Paths cache saved to", jsonFilePath)
}
