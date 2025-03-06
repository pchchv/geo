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
