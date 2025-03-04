package simplify

import "math"

// VisvalingamSimplifier is a reducer that performs the vivalingham algorithm.
type VisvalingamSimplifier struct {
	Threshold float64
	ToKeep    int // If 0, the default will be 2 for line, 3 for non-closed rings and 4 for closed rings.
	// The intent is to maintain valid geometry after simplification,
	// however it is still possible for the simplification to create self-intersections.
}

// Visvalingam creates a new VisvalingamSimplifier.
// If minPointsToKeep is 0 the algorithm will keep at least 2 points for lines,
// 3 for non-closed rings and 4 for closed rings.
// However it is still possible for the simplification to create self-intersections.
func Visvalingam(threshold float64, minPointsToKeep int) *VisvalingamSimplifier {
	return &VisvalingamSimplifier{
		Threshold: threshold,
		ToKeep:    minPointsToKeep,
	}
}

// VisvalingamThreshold runs the Visvalingam-Whyatt algorithm
// removing triangles whose area is below the threshold.
// Will keep at least 2 points for lines, 3 for non-closed rings and 4 for closed rings.
// The intent is to maintain valid geometry after simplification,
// however it is still possible for the simplification to create self-intersections.
func VisvalingamThreshold(threshold float64) *VisvalingamSimplifier {
	return Visvalingam(threshold, 0)
}

// VisvalingamKeep runs the Visvalingam-Wyatt algorithm,
// removing triangles of minimum area until the number of points `minPointsToKeep` is reached.
// If minPointsToKeep is 0, the algorithm will keep at least 2 points for lines,
// 3 for unclosed rings, and 4 for closed rings.
// However, it is still possible to create self-intersections when simplifying.
func VisvalingamKeep(minPointsToKeep int) *VisvalingamSimplifier {
	return Visvalingam(math.MaxFloat64, minPointsToKeep)
}

type visItem struct {
	area       float64  // triangle area
	pointIndex int      // index of point in original path
	next       *visItem // to keep a virtual linked list to help restore the triangle areas when points are removed
	previous   *visItem
	index      int // internal index in heap, for removal and update
}

// minHeap creates a priority queue or min heap.
type minHeap []*visItem
