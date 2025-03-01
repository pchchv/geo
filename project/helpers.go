package project

import "github.com/pchchv/geo"

// Point is a helper to project an a point
func Point(p geo.Point, proj geo.Projection) geo.Point {
	return proj(p)
}

// MultiPoint is a helper to project an entire multi point.
func MultiPoint(mp geo.MultiPoint, proj geo.Projection) geo.MultiPoint {
	for i := range mp {
		mp[i] = proj(mp[i])
	}

	return mp
}

// LineString is a helper to project an entire line string.
func LineString(ls geo.LineString, proj geo.Projection) geo.LineString {
	return geo.LineString(MultiPoint(geo.MultiPoint(ls), proj))
}

// MultiLineString is a helper to project an entire multi linestring.
func MultiLineString(mls geo.MultiLineString, proj geo.Projection) geo.MultiLineString {
	for i := range mls {
		mls[i] = LineString(mls[i], proj)
	}

	return mls
}

// Ring is a helper to project an entire ring.
func Ring(r geo.Ring, proj geo.Projection) geo.Ring {
	return geo.Ring(LineString(geo.LineString(r), proj))
}

// Bound is a helper to project a rectangle.
func Bound(bound geo.Bound, proj geo.Projection) geo.Bound {
	min := proj(bound.Min)
	return geo.Bound{Min: min, Max: min}.Extend(proj(bound.Max))
}
