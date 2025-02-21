package geo

// MultiLineString is a set of polylines.
type MultiLineString []LineString

// Bound returns a bound around all the line strings.
func (mls MultiLineString) Bound() Bound {
	if len(mls) == 0 {
		return emptyBound
	}

	bound := mls[0].Bound()
	for i := 1; i < len(mls); i++ {
		bound = bound.Union(mls[i].Bound())
	}

	return bound
}

// Equal compares two multi line strings.
// Returns true if lengths are the same and all points are Equal.
func (mls MultiLineString) Equal(multiLineString MultiLineString) bool {
	if len(mls) != len(multiLineString) {
		return false
	}

	for i, ls := range mls {
		if !ls.Equal(multiLineString[i]) {
			return false
		}
	}

	return true
}
