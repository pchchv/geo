package clip

import (
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
