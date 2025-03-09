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

// Midpoint returns the half-way point along a great circle path between the two points.
func Midpoint(p, p2 geo.Point) geo.Point {
	dLon := deg2rad(p2[0] - p[0])
	aLatRad := deg2rad(p[1])
	bLatRad := deg2rad(p2[1])
	x := math.Cos(bLatRad) * math.Cos(dLon)
	y := math.Cos(bLatRad) * math.Sin(dLon)
	r := geo.Point{
		deg2rad(p[0]) + math.Atan2(y, math.Cos(aLatRad)+x),
		math.Atan2(math.Sin(aLatRad)+math.Sin(bLatRad), math.Sqrt((math.Cos(aLatRad)+x)*(math.Cos(aLatRad)+x)+y*y)),
	}

	// convert back to degrees
	r[0] = rad2deg(r[0])
	r[1] = rad2deg(r[1])

	return r
}

// PointAtBearingAndDistance returns the point at the given bearing and distance in meters from the point
func PointAtBearingAndDistance(p geo.Point, bearing, distance float64) geo.Point {
	aLat := deg2rad(p[1])
	aLon := deg2rad(p[0])
	bearingRadians := deg2rad(bearing)
	distanceRatio := distance / geo.EarthRadius
	bLat := math.Asin(math.Sin(aLat)*math.Cos(distanceRatio) + math.Cos(aLat)*math.Sin(distanceRatio)*math.Cos(bearingRadians))
	bLon := aLon + math.Atan2(math.Sin(bearingRadians)*math.Sin(distanceRatio)*math.Cos(aLat), math.Cos(distanceRatio)-math.Sin(aLat)*math.Sin(bLat))
	return geo.Point{rad2deg(bLon), rad2deg(bLat)}
}
