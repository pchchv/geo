package geo

// Ring represents a set of ring on the earth.
type Ring LineString

// Bound returns a rect around the ring.
// Uses rectangular coordinates.
func (r Ring) Bound() Bound {
	return MultiPoint(r).Bound()
}

// Equal compares two rings.
// Returns true if lengths are the same and all points are Equal.
func (r Ring) Equal(ring Ring) bool {
	return MultiPoint(r).Equal(MultiPoint(ring))
}

// Dimensions returns 2 because a Ring is a 2d object.
func (r Ring) Dimensions() int {
	return 2
}

// GeoJSONType returns the GeoJSON type for the object.
func (r Ring) GeoJSONType() string {
	return "Polygon"
}
