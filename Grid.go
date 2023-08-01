package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/beefsack/go-astar"
)

// Tile represents a single tile with various properties.
type Tile struct {
	Walkable bool
	Snipe    bool
	Water    bool
	Cliff    bool
	x, y     int
	w        *Grid
}

// Grid represents the game map as a 2D grid of tiles.
type Grid struct {
	Width  int
	Height int
	Data   []Tile
}

// returns a tile for x,y
func (g *Grid) Tile(x, y int) *Tile {
	if x < 0 || y < 0 || x >= int(g.Width) || y >= int(g.Height) {
		return &g.Data[y*g.Width+x]
	}
	return nil
}

func (t *Tile) PathNeighbors() []astar.Pather {
	homies := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.w.Tile(t.x+offset[0], t.y+offset[1]); n != nil && n.Walkable {
			homies = append(homies, n)
		}
	}
	return homies
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toTile := to.(*Tile)
	return float64((toTile.x-t.x)*(toTile.x-t.x) + (toTile.y-t.y)*(toTile.y-t.y))
}

// NewGrid creates a new grid with the given dimensions.
func NewGrid(width, height int) *Grid {
	grid := &Grid{
		Width:  width,
		Height: height,
		Data:   make([]Tile, width*height),
	}
	return grid
}

// SetTileProperties sets the properties of a tile at the given coordinates.
func (grid *Grid) SetTileProperties(index int, walkable, snipe, water, cliff bool) {
	if index >= 0 && index < len(grid.Data) {
		grid.Data[index].Walkable = walkable
		grid.Data[index].Snipe = snipe
		grid.Data[index].Water = water
		grid.Data[index].Cliff = cliff
		grid.Data[index].x = index % grid.Width
		grid.Data[index].y = index / grid.Width
		grid.Data[index].w = grid
	}
}

// GetTileProperties returns the properties of a tile at the given coordinates.
func (grid *Grid) GetTileProperties(x, y int) (walkable, snipe, water, cliff bool) {
	if x >= 0 && y >= 0 && x < int(grid.Width) && y < int(grid.Height) {
		tile := grid.Data[y*grid.Width+x]
		// tile := grid.Data[y][x]
		return tile.Walkable, tile.Snipe, tile.Water, tile.Cliff
	}
	return false, false, false, false
}

// ParseGridFromFile parses the ".fld2" file and constructs the grid.
func ParseGridFromFile(filePath string) (*Grid, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("file is too short to contain width and height information")
	}

	// Extract width and height from the first four bytes.
	widthData := []byte{data[1], data[0]}
	width := int(binary.BigEndian.Uint16(widthData))
	heightData := []byte{data[3], data[2]}
	height := int(binary.BigEndian.Uint16(heightData))
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	if len(data) != int(4+width*height) {
		return nil, fmt.Errorf("invalid data format. The file size doesn't match width*height")
	}

	grid := NewGrid(width, height)
	parsedData := data[4:]

	for i := 0; i < len(parsedData); i++ {
		tileByte := parsedData[i]
		walkable := tileByte&TILE_WALK > 0
		snipe := tileByte&TILE_SNIPE > 0
		water := tileByte&TILE_WATER > 0
		cliff := tileByte&TILE_CLIFF > 0
		grid.SetTileProperties(i, walkable, snipe, water, cliff)
	}

	return grid, nil
}

// Constants for tile properties.
const (
	TILE_NOWALK = 0
	TILE_WALK   = 1 << 0
	TILE_SNIPE  = 1 << 1
	TILE_WATER  = 1 << 2
	TILE_CLIFF  = 1 << 3
)

func pathToVectors(path []astar.Pather) []*Vector {
	vectors := make([]*Vector, len(path))
	for i, node := range path {
		idk := node.(*Tile)
		vectors[i] = vectorPool.Get(idk.x, idk.y)
	}
	return vectors
}

func Idk(grid *Grid, start, end *Vector) ([]*Vector, float64, bool) {
	path, distance, found := astar.Path(grid.Tile(start.X, start.Y), grid.Tile(end.X, end.Y))
	if !found {
		return []*Vector{}, 0, false
	}

	return pathToVectors(path), distance, found
}
