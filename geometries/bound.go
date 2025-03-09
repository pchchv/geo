package geometries

import (
	"math"

	"github.com/pchchv/geo"
)

var (
	minLatitude  = deg2rad(-90)
	maxLatitude  = deg2rad(90)
	minLongitude = deg2rad(-180)
	maxLongitude = deg2rad(180)
)

// NewBoundAroundPoint creates a new bound given a center point,
// and a distance from the center point in meters.
func NewBoundAroundPoint(center geo.Point, distance float64) geo.Bound {
	var minLon, maxLon float64
	radDist := distance / geo.EarthRadius
	radLat := deg2rad(center[1])
	radLon := deg2rad(center[0])
	minLat := radLat - radDist
	maxLat := radLat + radDist
	if minLat > minLatitude && maxLat < maxLatitude {
		deltaLon := math.Asin(math.Sin(radDist) / math.Cos(radLat))
		minLon = radLon - deltaLon
		if minLon < minLongitude {
			minLon += 2 * math.Pi
		}

		maxLon = radLon + deltaLon
		if maxLon > maxLongitude {
			maxLon -= 2 * math.Pi
		}
	} else {
		minLat = math.Max(minLat, minLatitude)
		maxLat = math.Min(maxLat, maxLatitude)
		minLon = minLongitude
		maxLon = maxLongitude
	}

	return geo.Bound{
		Min: geo.Point{rad2deg(minLon), rad2deg(minLat)},
		Max: geo.Point{rad2deg(maxLon), rad2deg(maxLat)},
	}
}

// BoundHeight returns the approximate height in meters.
func BoundHeight(b geo.Bound) float64 {
	return 111131.75 * (b.Max[1] - b.Min[1])
}

// BoundWidth returns the approximate width in meters
// of the center of the bound.
func BoundWidth(b geo.Bound) float64 {
	c := (b.Min[1] + b.Max[1]) / 2.0
	s1 := geo.Point{b.Min[0], c}
	s2 := geo.Point{b.Max[0], c}
	return Distance(s1, s2)
}

// BoundPad expands the bound in all directions by the given amount of meters.
func BoundPad(b geo.Bound, meters float64) geo.Bound {
	dy := meters / 111131.75
	dx := dy / math.Cos(deg2rad(b.Max[1]))
	dx = math.Max(dx, dy/math.Cos(deg2rad(b.Min[1])))
	b.Min[0] -= dx
	b.Min[1] -= dy
	b.Max[0] += dx
	b.Max[1] += dy
	b.Min[0] = math.Max(b.Min[0], -180)
	b.Min[1] = math.Max(b.Min[1], -90)
	b.Max[0] = math.Min(b.Max[0], 180)
	b.Max[1] = math.Min(b.Max[1], 90)
	return b
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(r float64) float64 {
	return 180.0 * r / math.Pi
}
