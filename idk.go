package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Path represents a single path from one walkable cell to another.
// type Path struct {
// 	StartX, StartY int
// 	GoalX, GoalY   int
// 	Nodes          []*Node
// }

// PathsCache is a map to store precomputed paths from each walkable cell to each other walkable cell.
type PathsCache map[string]map[string][]*Vector

func reverseArray(a []*Vector) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func getKey(x, y int) string {
	return fmt.Sprintf("(%d,%d)", x, y)
}

func traversePathForward(path []*Vector, endX, endY int, paths PathsCache) {
	endKey := getKey(endX, endY)

	for i := 0; i < len(path); i++ {
		startKey := getKey(path[i].X, path[i].Y)
		if _, ok := paths[startKey]; !ok {
			paths[startKey] = make(map[string][]*Vector)
		}
		if _, ok := paths[startKey][endKey]; !ok {
			paths[startKey][endKey] = path[i:]
		}
	}
}

func traversePathBackward(path []*Vector, startX, startY int, paths PathsCache) {
	startKey := getKey(startX, startY)
	if _, ok := paths[startKey]; !ok {
		paths[startKey] = make(map[string][]*Vector)
	}
	for i := len(path) - 1; i >= 0; i-- {
		endKey := getKey(path[i].X, path[i].Y)
		if _, ok := paths[startKey][endKey]; !ok {
			paths[startKey][endKey] = path[:i+1]
		}
	}
}

func solvePath(grid *Grid, startX, startY, endX, endY int, paths PathsCache) {
	startKey := getKey(startX, startY)
	goalKey := getKey(endX, endY)

	// Check if the start key has already been initialized.
	if _, ok := paths[startKey]; !ok {
		paths[startKey] = make(map[string][]*Vector)
	}

	// Check if the path has already been calculated and stored.
	if _, ok := paths[startKey][goalKey]; !ok {
		vecPath, _, found := Idk(grid, &Vector{startX, startY}, &Vector{endX, endY})
		if !found {
			fmt.Printf("Path not found from (%d,%d) to (%d,%d)", startX, startY, endX, endY)
			return
		}
		paths[startKey][goalKey] = vecPath
		traversePathForward(vecPath, endX, endY, paths)
		traversePathBackward(vecPath, startX, startY, paths)
		reverseArray(vecPath)
		traversePathForward(vecPath, startX, startY, paths)
		traversePathBackward(vecPath, endX, endY, paths)
	}
}

// vector pool

// PrecomputePaths calculates and stores all paths in the PathsCache.
func PrecomputePaths(grid *Grid) PathsCache {
	vectorPool = NewVectorPool(grid.Width, grid.Height)
	paths := make(PathsCache)

	for y1 := 0; y1 < int(grid.Height); y1++ {
		for x1 := 0; x1 < int(grid.Width); x1++ {
			if !grid.Tile(x1, y1).Walkable {
				continue
			}

			for y2 := 0; y2 < int(grid.Height); y2++ {
				for x2 := 0; x2 < int(grid.Width); x2++ {
					if !grid.Tile(x2, y2).Walkable {
						continue
					}
					solvePath(grid, x1, y1, x2, y2, paths)

				}
			}
		}
	}

	return paths
}

// SavePathsCacheToJSON dumps the PathsCache to a JSON file.
func SavePathsCacheToJSON(cache PathsCache, filePath string) error {
	jsonData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadPathsCacheFromJSON loads the PathsCache from a JSON file.
func LoadPathsCacheFromJSON(filePath string) (PathsCache, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cache := make(PathsCache)
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
