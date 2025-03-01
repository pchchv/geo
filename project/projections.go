package project

import (
	"math"

	"github.com/pchchv/geo"
)

const earthRadiusPi = geo.EarthRadius * math.Pi

// Mercator performs the Spherical Pseudo-Mercator projection used by most web maps.
var Mercator = struct {
	ToWGS84 geo.Projection
}{
	ToWGS84: func(p geo.Point) geo.Point {
		return geo.Point{
			180.0 * p[0] / earthRadiusPi,
			180.0 / math.Pi * (2*math.Atan(math.Exp(p[1]/geo.EarthRadius)) - math.Pi/2.0),
		}
	},
}
