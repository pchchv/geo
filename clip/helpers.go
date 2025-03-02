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
