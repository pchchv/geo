package planar

import (
	"math"

	"github.com/pchchv/geo"
)

// DistanceFromSegmentSquared returns point's squared distance from the segement [a, b].
func DistanceFromSegmentSquared(a, b, point geo.Point) float64 {
	return segmentDistanceFromSquared(a, b, point)
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

func multiPointDistanceFrom(mp geo.MultiPoint, p geo.Point) (float64, int) {
	index := -1
	dist := math.Inf(1)
	for i := range mp {
		if d := DistanceSquared(mp[i], p); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}

func lineStringDistanceFrom(ls geo.LineString, p geo.Point) (float64, int) {
	index := -1
	dist := math.Inf(1)
	for i := 0; i < len(ls)-1; i++ {
		if d := segmentDistanceFromSquared(ls[i], ls[i+1], p); d < dist {
			dist = d
			index = i
		}
	}

	return math.Sqrt(dist), index
}
