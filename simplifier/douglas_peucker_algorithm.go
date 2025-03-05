package simplifier

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

var _ geo.Simplifier = &DouglasPeuckerSimplifier{}

// DouglasPeuckerSimplifier wraps the DouglasPeucker function.
type DouglasPeuckerSimplifier struct {
	Threshold float64
}

// DouglasPeucker creates a new DouglasPeuckerSimplifier.
func DouglasPeucker(threshold float64) *DouglasPeuckerSimplifier {
	return &DouglasPeuckerSimplifier{
		Threshold: threshold,
	}
}

// LineString will simplify the linestring using this simplifier.
func (s *DouglasPeuckerSimplifier) LineString(ls geo.LineString) geo.LineString {
	return lineString(s, ls)
}

// MultiLineString will simplify the multi-linestring using this simplifier.
func (s *DouglasPeuckerSimplifier) MultiLineString(mls geo.MultiLineString) geo.MultiLineString {
	return multiLineString(s, mls)
}

// Polygon will simplify the polygon using this simplifier.
func (s *DouglasPeuckerSimplifier) Polygon(p geo.Polygon) geo.Polygon {
	return polygon(s, p)
}

// MultiPolygon will simplify the multi-polygon using this simplifier.
func (s *DouglasPeuckerSimplifier) MultiPolygon(mp geo.MultiPolygon) geo.MultiPolygon {
	return multiPolygon(s, mp)
}

// Ring will simplify the ring using this simplifier.
func (s *DouglasPeuckerSimplifier) Ring(r geo.Ring) geo.Ring {
	return ring(s, r)
}

// Collection will simplify the collection using this simplifier.
func (s *DouglasPeuckerSimplifier) Collection(c geo.Collection) geo.Collection {
	return collection(s, c)
}

// Simplify will run the simplification for any geometry type.
func (s *DouglasPeuckerSimplifier) Simplify(g geo.Geometry) geo.Geometry {
	return simplify(s, g)
}

func (s *DouglasPeuckerSimplifier) simplify(ls geo.LineString, area, wim bool) (geo.LineString, []int) {
	var indexMap []int
	mask := make([]byte, len(ls))
	mask[0] = 1
	mask[len(mask)-1] = 1
	found := dpWorker(ls, s.Threshold, mask)
	if wim {
		indexMap = make([]int, 0, found)
	}

	var count int
	for i, v := range mask {
		if v == 1 {
			ls[count] = ls[i]
			count++
			if wim {
				indexMap = append(indexMap, i)
			}
		}
	}

	return ls[:count], indexMap
}

// dpWorker performs recursive threshold checks.
func dpWorker(ls geo.LineString, threshold float64, mask []byte) int {
	var stack []int
	found := 2
	stack = append(stack, 0, len(ls)-1)
	for len(stack) > 0 {
		var maxIndex int
		var maxDist float64
		start := stack[len(stack)-2]
		end := stack[len(stack)-1]
		for i := start + 1; i < end; i++ {
			dist := planar.DistanceFromSegmentSquared(ls[start], ls[end], ls[i])
			if dist > maxDist {
				maxDist = dist
				maxIndex = i
			}
		}

		if maxDist > threshold*threshold {
			found++
			mask[maxIndex] = 1
			stack[len(stack)-1] = maxIndex
			stack = append(stack, maxIndex, end)
		} else {
			stack = stack[:len(stack)-2]
		}
	}

	return found
}
