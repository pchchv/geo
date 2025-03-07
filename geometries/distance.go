package geometries

import (
	"math"

	"github.com/pchchv/geo"
)

// Distance returns the distance between two points on the earth.
func Distance(p1, p2 geo.Point) float64 {
	dLon := math.Abs(deg2rad(p1[0] - p2[0]))
	if dLon > math.Pi {
		dLon = 2*math.Pi - dLon
	}

	dLat := deg2rad(p1[1] - p2[1])
	// fast way using pythagorean theorem on an equirectangular projection
	x := dLon * math.Cos(deg2rad((p1[1]+p2[1])/2.0))
	return math.Sqrt(dLat*dLat+x*x) * geo.EarthRadius
}

// DistanceHaversine computes the distance on the earth using the
// more accurate haversine formula.
func DistanceHaversine(p1, p2 geo.Point) float64 {
	dLat := deg2rad(p1[1] - p2[1])
	dLon := deg2rad(p1[0] - p2[0])
	dLat2Sin := math.Sin(dLat / 2)
	dLon2Sin := math.Sin(dLon / 2)
	a := dLat2Sin*dLat2Sin + math.Cos(deg2rad(p2[1]))*math.Cos(deg2rad(p1[1]))*dLon2Sin*dLon2Sin

	return 2.0 * geo.EarthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
