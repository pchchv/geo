package clip

import (
	"fmt"
	"math"

	"github.com/pchchv/geo"
)

// Ring clips the ring to the bounding box and returns another ring.
// This operation will modify the input by
// using as a scratch space so clone if necessary.
func Ring(b geo.Bound, r geo.Ring) geo.Ring {
	result := ring(b, r)
	if len(result) == 0 {
		return nil
	}

	return result
}

// LineString clips the linestring to the bounding box.
func LineString(b geo.Bound, ls geo.LineString, opts ...Option) geo.MultiLineString {
	var open bool
	if len(opts) > 0 {
		o := &options{}
		for _, opt := range opts {
			opt(o)
		}

		open = o.openBound
	}

	result := line(b, ls, open)
	if len(result) == 0 {
		return nil
	}

	return result
}

// MultiLineString clips the linestrings to the bounding box and returns a linestring union.
func MultiLineString(b geo.Bound, mls geo.MultiLineString, opts ...Option) (result geo.MultiLineString) {
	var open bool
	if len(opts) > 0 {
		o := &options{}
		for _, opt := range opts {
			opt(o)
		}

		open = o.openBound
	}

	for _, ls := range mls {
		r := line(b, ls, open)
		if len(r) != 0 {
			result = append(result, r...)
		}
	}

	return
}

// Polygon clips the polygon to the
// bounding box excluding the inner rings if they
// do not intersect the bounding box.
// This operation will modify the input by
// using as a scratch space so clone if necessary.
func Polygon(b geo.Bound, p geo.Polygon) geo.Polygon {
	if len(p) == 0 {
		return nil
	}

	r := Ring(b, p[0])
	if r == nil {
		return nil
	}

	result := geo.Polygon{r}
	for i := 1; i < len(p); i++ {
		if r = Ring(b, p[i]); r != nil {
			result = append(result, r)
		}
	}

	return result
}

// MultiPolygon clips the multi polygon to the bounding box excluding
// any polygons if they don't intersect the bounding box.
// This operation will modify the input by
// using as a scratch space so clone if necessary.
func MultiPolygon(b geo.Bound, mp geo.MultiPolygon) (result geo.MultiPolygon) {
	for _, polygon := range mp {
		if p := Polygon(b, polygon); p != nil {
			result = append(result, p)
		}
	}

	return
}

// MultiPoint returns a new set with the points outside the bound removed.
func MultiPoint(b geo.Bound, mp geo.MultiPoint) (result geo.MultiPoint) {
	for _, p := range mp {
		if b.Contains(p) {
			result = append(result, p)
		}
	}

	return
}

// Bound intersects the two bounds.
// May result in an empty/degenerate bound.
func Bound(b, bound geo.Bound) geo.Bound {
	if b.IsEmpty() {
		return bound
	} else if bound.IsEmpty() {
		return b
	}

	return geo.Bound{
		Min: geo.Point{
			math.Max(b.Min[0], bound.Min[0]),
			math.Max(b.Min[1], bound.Min[1]),
		},
		Max: geo.Point{
			math.Min(b.Max[0], bound.Max[0]),
			math.Min(b.Max[1], bound.Max[1]),
		},
	}
}

// Collection clips each element in the collection to the bounding box.
// It will exclude elements if they don't intersect the bounding box.
// This operation will modify the input of '2d geometry' by
// using as a scratch space so clone if necessary.
func Collection(b geo.Bound, c geo.Collection) (result geo.Collection) {
	for _, g := range c {
		if clipped := Geometry(b, g); clipped != nil {
			result = append(result, clipped)
		}
	}

	return
}

// Geometry will clip the geometry to the bounding box
// using the correct functions for the type.
// This operation will modify the input of '1d or 2d geometry' by
// using as a scratch space so clone if necessary.
func Geometry(b geo.Bound, g geo.Geometry) geo.Geometry {
	if g == nil || !b.Intersects(g.Bound()) {
		return nil
	}

	switch g := g.(type) {
	case geo.Point:
		return g // intersect check above
	case geo.MultiPoint:
		if mp := MultiPoint(b, g); len(mp) == 1 {
			return mp[0]
		} else if mp == nil {
			return nil
		} else {
			return mp
		}
	case geo.LineString:
		if mls := LineString(b, g); len(mls) == 1 {
			return mls[0]
		} else if len(mls) == 0 {
			return nil
		} else {
			return mls
		}
	case geo.MultiLineString:
		if mls := MultiLineString(b, g); len(mls) == 1 {
			return mls[0]
		} else if mls == nil {
			return nil
		} else {
			return mls
		}
	case geo.Ring:
		return Ring(b, g)
	case geo.Polygon:
		return Polygon(b, g)
	case geo.MultiPolygon:
		if mp := MultiPolygon(b, g); len(mp) == 1 {
			return mp[0]
		} else if mp == nil {
			return nil
		} else {
			return mp
		}
	case geo.Collection:
		if c := Collection(b, g); len(c) == 1 {
			return c[0]
		} else if c == nil {
			return nil
		} else {
			return c
		}
	case geo.Bound:
		return Bound(b, g)
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}
