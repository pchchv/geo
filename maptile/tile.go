package maptile

import (
	"math"
	"math/bits"

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

// Parent returns the parent of the tile.
func (t Tile) Parent() Tile {
	if t.Z == 0 {
		return t
	}

	return Tile{
		X: t.X >> 1,
		Y: t.Y >> 1,
		Z: t.Z - 1,
	}
}

// Center returns the center of the tile.
func (t Tile) Center() geo.Point {
	return t.Bound(0).Center()
}

// Valid returns if the tile's x/y are
// within the range for the tile's zoom.
func (t Tile) Valid() bool {
	maxIndex := uint32(1) << uint32(t.Z)
	return t.X < maxIndex && t.Y < maxIndex
}

// Quadkey returns the quad key for the tile.
func (t Tile) Quadkey() uint64 {
	var i, result uint64
	for i = 0; i < uint64(t.Z); i++ {
		result |= (uint64(t.X) & (1 << i)) << i
		result |= (uint64(t.Y) & (1 << i)) << (i + 1)
	}

	return result
}

// Contains returns if the given tile is
// fully contained (or equal to) the give tile.
func (t Tile) Contains(tile Tile) bool {
	return tile.Z >= t.Z && t == tile.toZoom(t.Z)
}

// Range returns the min and max tile "range" to
// cover the tile at the given zoom.
func (t Tile) Range(z Zoom) (min, max Tile) {
	if z < t.Z {
		t = t.toZoom(z)
		return t, t
	}

	offset := z - t.Z
	return Tile{
			X: t.X << offset,
			Y: t.Y << offset,
			Z: z,
		}, Tile{
			X: ((t.X + 1) << offset) - 1,
			Y: ((t.Y + 1) << offset) - 1,
			Z: z,
		}
}

// SharedParent returns the tile that contains both the tiles.
func (t Tile) SharedParent(tile Tile) Tile {
	// bring both tiles to the lowest zoom.
	if t.Z != tile.Z {
		if t.Z < tile.Z {
			tile = tile.toZoom(t.Z)
		} else {
			t = t.toZoom(tile.Z)
		}
	}

	if t == tile {
		return t
	}

	// bits different for x and y
	xc := uint32(32 - bits.LeadingZeros32(t.X^tile.X))
	yc := uint32(32 - bits.LeadingZeros32(t.Y^tile.Y))

	// max of xc, yc
	maxc := xc
	if yc > maxc {
		maxc = yc
	}

	return Tile{
		X: t.X >> maxc,
		Y: t.Y >> maxc,
		Z: t.Z - Zoom(maxc),
	}
}

func (t Tile) toZoom(z Zoom) Tile {
	if z > t.Z {
		return Tile{
			X: t.X << (z - t.Z),
			Y: t.Y << (z - t.Z),
			Z: z,
		}
	}

	return Tile{
		X: t.X >> (t.Z - z),
		Y: t.Y >> (t.Z - z),
		Z: z,
	}
}

// At creates a tile for the point at the given zoom.
// Will create a valid tile for the zoom.
// Points outside the range lat [-85.0511, 85.0511]
// will be snapped to the max or min tile as appropriate.
func At(ll geo.Point, z Zoom) Tile {
	f := Fraction(ll, z)
	t := Tile{
		X: uint32(f[0]),
		Y: uint32(f[1]),
		Z: z,
	}

	return t
}

// Fraction returns the precise tile fraction at the given zoom.
// Will return 2^zoom-1 if the point is below 85.0511 S.
func Fraction(ll geo.Point, z Zoom) (p geo.Point) {
	factor := uint32(1 << z)
	maxtiles := float64(factor)
	lng := ll[0]/360.0 + 0.5
	p[0] = lng * maxtiles
	// bound it because we have a top of the world problem
	if ll[1] < -85.0511 {
		p[1] = maxtiles - 1
	} else if ll[1] > 85.0511 {
		p[1] = 0
	} else {
		siny := math.Sin(ll[1] * math.Pi / 180.0)
		lat := 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		p[1] = lat * maxtiles
	}

	return
}

// FromQuadkey creates the tile from the quadkey.
func FromQuadkey(k uint64, z Zoom) Tile {
	t := Tile{Z: z}
	for i := Zoom(0); i < z; i++ {
		t.X |= uint32((k & (1 << (2 * i))) >> i)
		t.Y |= uint32((k & (1 << (2*i + 1))) >> (i + 1))
	}

	return t
}
