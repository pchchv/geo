package geometries

import (
	"fmt"
	"math"

	"github.com/pchchv/geo"
)

// SignedArea returns the signed area of the ring.
// Return negative if the ring is in the clockwise direction.
// Implicitly close the ring.
func SignedArea(r geo.Ring) float64 {
	return ringArea(r)
}

// Area returns the area of the geometry on the earth.
func Area(g geo.Geometry) float64 {
	switch g := g.(type) {
	case nil, geo.Point, geo.MultiPoint, geo.LineString, geo.MultiLineString:
		return 0
	case geo.Ring:
		return math.Abs(ringArea(g))
	case geo.Polygon:
		return polygonArea(g)
	case geo.MultiPolygon:
		return multiPolygonArea(g)
	case geo.Collection:
		return collectionArea(g)
	case geo.Bound:
		return Area(g.ToRing())
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}

func ringArea(r geo.Ring) float64 {
	if len(r) < 3 {
		return 0
	}

	var lo, mi, hi int
	l := len(r)
	if r[0] != r[len(r)-1] {
		// if not a closed ring,
		// add an implicit calc for that last point
		l++
	}

	// to support implicit closing of ring,
	// replace references to the last point in r to the first 1
	var area float64
	for i := 0; i < l; i++ {
		if i == l-3 { // i = N-3
			lo = l - 3
			mi = l - 2
			hi = 0
		} else if i == l-2 { // i = N-2
			lo = l - 2
			mi = 0
			hi = 0
		} else if i == l-1 { // i = N-1
			lo = 0
			mi = 0
			hi = 1
		} else { // i = 0 to N-3
			lo = i
			mi = i + 1
			hi = i + 2
		}

		area += (deg2rad(r[hi][0]) - deg2rad(r[lo][0])) * math.Sin(deg2rad(r[mi][1]))
	}

	return -area * geo.EarthRadius * geo.EarthRadius / 2
}

func collectionArea(c geo.Collection) (area float64) {
	for _, g := range c {
		area += Area(g)
	}
	return
}

func polygonArea(p geo.Polygon) float64 {
	if len(p) == 0 {
		return 0
	}

	sum := math.Abs(ringArea(p[0]))
	for i := 1; i < len(p); i++ {
		sum -= math.Abs(ringArea(p[i]))
	}

	return sum
}

func multiPolygonArea(mp geo.MultiPolygon) (sum float64) {
	for _, p := range mp {
		sum += polygonArea(p)
	}
	return
}
