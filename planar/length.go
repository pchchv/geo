package planar

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/internal/length"
)

// Length returns the length of the boundary of
// the geometry using 2d euclidean geometry.
func Length(g geo.Geometry) float64 {
	return length.Length(g, Distance)
}
