package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Tile represents a single tile with various properties.
type Tile struct {
	Walkable bool
	Snipe    bool
	Water    bool
	Cliff    bool
}

// Grid represents the game map as a 2D grid of tiles.
type Grid struct {
	Width  int
	Height int
	Data   [][]Tile
}

// NewGrid creates a new grid with the given dimensions.
func NewGrid(width, height int) *Grid {
	grid := &Grid{
		Width:  width,
		Height: height,
		Data:   make([][]Tile, height),
	}
	for i := range grid.Data {
		grid.Data[i] = make([]Tile, width)
	}
	return grid
}

// SetTileProperties sets the properties of a tile at the given coordinates.
func (grid *Grid) SetTileProperties(x, y int, walkable, snipe, water, cliff bool) {
	if x >= 0 && y >= 0 && x < int(grid.Width) && y < int(grid.Height) {
		grid.Data[y][x].Walkable = walkable
		grid.Data[y][x].Snipe = snipe
		grid.Data[y][x].Water = water
		grid.Data[y][x].Cliff = cliff
	}
}

// GetTileProperties returns the properties of a tile at the given coordinates.
func (grid *Grid) GetTileProperties(x, y int) (walkable, snipe, water, cliff bool) {
	if x >= 0 && y >= 0 && x < int(grid.Width) && y < int(grid.Height) {
		tile := grid.Data[y][x]
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

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			tileByte := parsedData[y*int(width)+x]
			walkable := tileByte&TILE_WALK > 0
			snipe := tileByte&TILE_SNIPE > 0
			water := tileByte&TILE_WATER > 0
			cliff := tileByte&TILE_CLIFF > 0
			grid.SetTileProperties(x, y, walkable, snipe, water, cliff)
		}
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
