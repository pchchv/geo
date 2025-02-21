package geo

// MultiPoint represents a set of points in
// the 2D Eucledian or Cartesian plane.
type MultiPoint []Point

// Bound returns a bound around the points.
// Uses rectangular coordinates.
func (mp MultiPoint) Bound() Bound {
	if len(mp) == 0 {
		return emptyBound
	}

	b := Bound{mp[0], mp[0]}
	for _, p := range mp {
		b = b.Extend(p)
	}

	return b
}

// Equal compares two MultiPoint objects.
// Returns true if lengths are the same and all points are Equal,
// and in the same order.
func (mp MultiPoint) Equal(multiPoint MultiPoint) bool {
	if len(mp) != len(multiPoint) {
		return false
	}

	for i, point := range mp {
		if !point.Equal(multiPoint[i]) {
			return false
		}
	}

	return true
}

// GeoJSONType returns the GeoJSON type for the object.
func (mp MultiPoint) GeoJSONType() string {
	return "MultiPoint"
}

// Dimensions returns 0 because a MultiPoint is a 0d object.
func (mp MultiPoint) Dimensions() int {
	return 0
}
