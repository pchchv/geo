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

// Polygon is a helper to project an entire polygon.
func Polygon(p geo.Polygon, proj geo.Projection) geo.Polygon {
	for i := range p {
		p[i] = Ring(p[i], proj)
	}

	return p
}

// MultiPolygon is a helper to project an entire multi polygon.
func MultiPolygon(mp geo.MultiPolygon, proj geo.Projection) geo.MultiPolygon {
	for i := range mp {
		mp[i] = Polygon(mp[i], proj)
	}

	return mp
}

// Collection is a helper to project a rectangle.
func Collection(c geo.Collection, proj geo.Projection) geo.Collection {
	for i := range c {
		c[i] = Geometry(c[i], proj)
	}

	return c
}

// Geometry is a helper to project any geomtry.
func Geometry(g geo.Geometry, proj geo.Projection) geo.Geometry {
	if g == nil {
		return nil
	}

	switch g := g.(type) {
	case geo.Point:
		return Point(g, proj)
	case geo.MultiPoint:
		return MultiPoint(g, proj)
	case geo.LineString:
		return LineString(g, proj)
	case geo.MultiLineString:
		return MultiLineString(g, proj)
	case geo.Ring:
		return Ring(g, proj)
	case geo.Polygon:
		return Polygon(g, proj)
	case geo.MultiPolygon:
		return MultiPolygon(g, proj)
	case geo.Collection:
		return Collection(g, proj)
	case geo.Bound:
		return Bound(g, proj)
	}

	panic("geometry type not supported")
}
