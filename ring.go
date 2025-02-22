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

// Orientation returns 1 if the the ring is in couter-clockwise order,
// return -1 if the ring is the clockwise order and
// 0 if the ring is degenerate and had no area.
func (r Ring) Orientation() Orientation {
	// this is a fast planar area computation,
	// which is okay for this use
	// implicitly move everything to near the
	// origin to help with roundoff
	var area float64
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		area += (r[i][0]-offsetX)*(r[i+1][1]-offsetY) -
			(r[i+1][0]-offsetX)*(r[i][1]-offsetY)
	}

	if area > 0 {
		return CCW
	} else if area < 0 {
		return CW
	}

	// degenerate case, no area
	return 0
}
