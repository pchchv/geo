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

// LineString will simplify the linestring using this simplifier.
func (s *RadialSimplifier) LineString(ls geo.LineString) geo.LineString {
	return lineString(s, ls)
}

// MultiLineString will simplify the multi-linestring using this simplifier.
func (s *RadialSimplifier) MultiLineString(mls geo.MultiLineString) geo.MultiLineString {
	return multiLineString(s, mls)
}

func (s *RadialSimplifier) simplify(ls geo.LineString, area, wim bool) (geo.LineString, []int) {
	var indexMap []int
	if wim {
		indexMap = append(indexMap, 0)
	}

	var current int
	count := 1
	for i := 1; i < len(ls); i++ {
		if s.DistanceFunc(ls[current], ls[i]) > s.Threshold {
			current = i
			ls[count] = ls[i]
			count++
			if wim {
				indexMap = append(indexMap, current)
			}
		}
	}

	if current != len(ls)-1 {
		ls[count] = ls[len(ls)-1]
		count++
		if wim {
			indexMap = append(indexMap, len(ls)-1)
		}
	}

	return ls[:count], indexMap
}
