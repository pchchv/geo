package simplify

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

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
