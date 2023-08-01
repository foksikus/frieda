package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func measureExecutionTime(function func()) time.Duration {
	startTime := time.Now()
	function()
	endTime := time.Now()
	return endTime.Sub(startTime)
}

var vectorPool *VectorPool
var grid *Grid

// generate handler function
func GenerateHandler(g *Grid, vp *VectorPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		x1, err := strconv.Atoi(vars["x1"])
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}
		y1, err := strconv.Atoi(vars["y1"])
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}
		x2, err := strconv.Atoi(vars["x2"])
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}
		y2, err := strconv.Atoi(vars["y2"])
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}

		fmt.Println("Finding path from", x1, y1, "to", x2, y2)

		start := time.Now()
		path, dis, found := Idk(g, &Vector{x1, y1}, &Vector{x2, y2})
		fmt.Println("Time:", time.Since(start))
		if !found {
			fmt.Println("Path not found")
			return
		}

		fmt.Println("Path length:", len(path))
		fmt.Println("Distance:", dis)

		// Your logic to process x and y goes here...
		response := fmt.Sprintf("You requested /find/%d/%d/%d/%d", x1, y1, x2, y2)
		fmt.Fprintln(w, response)
	}
}

func main() {
	filePath := *flag.String("file", "prt_fild02.fld2", "filepath to the .fld2 file")

	grid, err := ParseGridFromFile(filePath)
	if err != nil {
		fmt.Println("Error parsing the file:", err)
		return
	}
	vectorPool = NewVectorPool(grid.Width, grid.Height)

	// fmt.Println("Finding path from", 355, 55, "to", 60, 320)

	// start := time.Now()
	// path, dis, found := Idk(grid, &Vector{355, 55}, &Vector{60, 320})
	// fmt.Println("Time:", time.Since(start))
	// if !found {
	// 	fmt.Println("Path not found")
	// 	return
	// }

	// fmt.Println("Path length:", len(path))
	// fmt.Println("Distance:", dis)

	r := mux.NewRouter()
	r.HandleFunc("/find/{x1}/{y1}/{x2}/{y2}", GenerateHandler(grid, vectorPool))

	http.Handle("/", r)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
