package simplify

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
