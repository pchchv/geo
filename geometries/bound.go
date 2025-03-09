package geometries

import (
	"math"

	"github.com/pchchv/geo"
)

// BoundHeight returns the approximate height in meters.
func BoundHeight(b geo.Bound) float64 {
	return 111131.75 * (b.Max[1] - b.Min[1])
}

// BoundWidth returns the approximate width in meters
// of the center of the bound.
func BoundWidth(b geo.Bound) float64 {
	c := (b.Min[1] + b.Max[1]) / 2.0
	s1 := geo.Point{b.Min[0], c}
	s2 := geo.Point{b.Max[0], c}
	return Distance(s1, s2)
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(r float64) float64 {
	return 180.0 * r / math.Pi
}
