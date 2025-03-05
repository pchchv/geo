package simplify

import "github.com/pchchv/geo"

// var _ geo.Simplifier = &RadialSimplifier{}

// RadialSimplifier wraps the Radial functions.
type RadialSimplifier struct {
	DistanceFunc geo.DistanceFunc
	Threshold    float64 // euclidean distance
}

// Radial creates a new RadialSimplifier.
func Radial(df geo.DistanceFunc, threshold float64) *RadialSimplifier {
	return &RadialSimplifier{
		DistanceFunc: df,
		Threshold:    threshold,
	}
}
