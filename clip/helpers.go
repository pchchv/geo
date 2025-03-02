package clip

import "github.com/pchchv/geo"

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
