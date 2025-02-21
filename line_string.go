package geo

// LineString represents a set of points to be thought of as a polyline.
type LineString []Point

// Bound returns a rect around the line string.
// Uses rectangular coordinates.
func (ls LineString) Bound() Bound {
	return MultiPoint(ls).Bound()
}

// Equal compares two line strings.
// Returns true if lengths are the same and all points are Equal.
func (ls LineString) Equal(lineString LineString) bool {
	return MultiPoint(ls).Equal(MultiPoint(lineString))
}

// GeoJSONType returns the GeoJSON type for the object.
func (ls LineString) GeoJSONType() string {
	return "LineString"
}

// Dimensions returns 1 because a LineString is a 1d object.
func (ls LineString) Dimensions() int {
	return 1
}
