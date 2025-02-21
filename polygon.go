package geo

// Polygon is a closed area.
// The first LineString is the outer ring.
// The others are the holes.
// Each LineString is expected to be closed
// ie. the first point matches the last.
type Polygon []Ring
