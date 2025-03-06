package maptile

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/internal/mercator"
)

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

// Bound returns the geo bound for the tile.
// An optional tileBuffer parameter can
// be passes to create a buffer around the bound in tile dimension.
// i.e. a tileBuffer of 1 would create a bound 9x the size of the tile,
// centered around the provided tile.
func (t Tile) Bound(tileBuffer ...float64) geo.Bound {
	var buffer float64
	if len(tileBuffer) > 0 {
		buffer = tileBuffer[0]
	}

	x := float64(t.X)
	y := float64(t.Y)
	minx := x - buffer
	miny := y - buffer
	if miny < 0 {
		miny = 0
	}

	lon1, lat1 := mercator.ToGeo(minx, miny, uint32(t.Z))
	maxx := x + 1 + buffer
	maxtiles := float64(uint32(1 << t.Z))
	maxy := y + 1 + buffer
	if maxy > maxtiles {
		maxy = maxtiles
	}

	lon2, lat2 := mercator.ToGeo(maxx, maxy, uint32(t.Z))
	return geo.Bound{
		Min: geo.Point{lon1, lat2},
		Max: geo.Point{lon2, lat1},
	}
}
