package simplifier

import "github.com/pchchv/geo"

var _ geo.Simplifier = &RadialSimplifier{}

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

// Polygon will simplify the polygon using this simplifier.
func (s *RadialSimplifier) Polygon(p geo.Polygon) geo.Polygon {
	return polygon(s, p)
}

// MultiPolygon will simplify the multi-polygon using this simplifier.
func (s *RadialSimplifier) MultiPolygon(mp geo.MultiPolygon) geo.MultiPolygon {
	return multiPolygon(s, mp)
}

// Ring will simplify the ring using this simplifier.
func (s *RadialSimplifier) Ring(r geo.Ring) geo.Ring {
	return ring(s, r)
}

// Collection will simplify the collection using this simplifier.
func (s *RadialSimplifier) Collection(c geo.Collection) geo.Collection {
	return collection(s, c)
}

// Simplify will run the simplification for any geometry type.
func (s *RadialSimplifier) Simplify(g geo.Geometry) geo.Geometry {
	return simplify(s, g)
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
