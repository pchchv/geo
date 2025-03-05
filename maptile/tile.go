package maptile

// Zoom is a strict type for a tile zoom level.
type Zoom uint32

// Tile is an x, y, z web mercator tile.
type Tile struct {
	X uint32
	Y uint32
	Z Zoom
}

// New creates a new tile with the given coordinates.
func New(x, y uint32, z Zoom) Tile {
	return Tile{x, y, z}
}
