package tilecover

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/maptile"
)

// Point creates a tile cover for the point,
// i.e. just the tile containing the point.
func Point(ll geo.Point, z maptile.Zoom) maptile.Set {
	return maptile.Set{
		maptile.At(ll, z): true,
	}
}

// MultiPoint creates a tile cover for the set of points,
func MultiPoint(mp geo.MultiPoint, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, p := range mp {
		set[maptile.At(p, z)] = true
	}

	return set
}

// Bound creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func Bound(b geo.Bound, z maptile.Zoom) maptile.Set {
	lo := maptile.At(b.Min, z)
	hi := maptile.At(b.Max, z)
	result := make(maptile.Set, (hi.X-lo.X+1)*(lo.Y-hi.Y+1))
	for x := lo.X; x <= hi.X; x++ {
		for y := hi.Y; y <= lo.Y; y++ {
			result[maptile.Tile{X: x, Y: y, Z: z}] = true
		}
	}

	return result
}
