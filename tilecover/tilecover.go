package tilecover

import (
	"fmt"

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

// Geometry returns the covering set of tiles for the given geometry.
func Geometry(g geo.Geometry, z maptile.Zoom) (maptile.Set, error) {
	switch g := g.(type) {
	case nil:
		return nil, nil
	case geo.Point:
		return Point(g, z), nil
	case geo.MultiPoint:
		return MultiPoint(g, z), nil
	case geo.LineString:
		return LineString(g, z), nil
	case geo.MultiLineString:
		return MultiLineString(g, z), nil
	case geo.Ring:
		return Ring(g, z)
	case geo.Polygon:
		return Polygon(g, z)
	case geo.MultiPolygon:
		return MultiPolygon(g, z)
	case geo.Collection:
		return Collection(g, z)
	case geo.Bound:
		return Bound(g, z), nil
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}

// Collection returns the covering set of tiles for the geometry collection.
func Collection(c geo.Collection, z maptile.Zoom) (maptile.Set, error) {
	set := make(maptile.Set)
	for _, g := range c {
		if s, err := Geometry(g, z); err != nil {
			return nil, err
		} else {
			set.Merge(s)
		}
	}

	return set, nil
}
