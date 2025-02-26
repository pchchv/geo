package planar

import (
	"math"

	"github.com/pchchv/geo"
)

// DistanceFromSegmentSquared returns point's squared distance from the segement [a, b].
func DistanceFromSegmentSquared(a, b, point geo.Point) float64 {
	x := a[0]
	y := a[1]
	dx := b[0] - x
	dy := b[1] - y
	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)
		if t > 1 {
			x = b[0]
			y = b[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y
	return dx*dx + dy*dy
}

// DistanceFromSegment returns the point's distance from the segment [a, b].
func DistanceFromSegment(a, b, point geo.Point) float64 {
	return math.Sqrt(DistanceFromSegmentSquared(a, b, point))
}

func segmentDistanceFromSquared(p1, p2, point geo.Point) float64 {
	x := p1[0]
	y := p1[1]
	dx := p2[0] - x
	dy := p2[1] - y
	if dx != 0 || dy != 0 {
		t := ((point[0]-x)*dx + (point[1]-y)*dy) / (dx*dx + dy*dy)
		if t > 1 {
			x = p2[0]
			y = p2[1]
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = point[0] - x
	dy = point[1] - y
	return dx*dx + dy*dy
}
