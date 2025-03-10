package geometries

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/internal/length"
)

// Length returns the length of the boundary of the geometry using the geo distance function.
func Length(g geo.Geometry) float64 {
	return length.Length(g, Distance)
}

// LengthHaversine returns the length of the boundary of the geometry using the geo haversine formula.
func LengthHaversine(g geo.Geometry) float64 {
	return length.Length(g, DistanceHaversine)
}
