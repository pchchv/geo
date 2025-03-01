package project

import (
	"fmt"
	"math"

	"github.com/pchchv/geo"
)

const earthRadiusPi = geo.EarthRadius * math.Pi

var (
	// Mercator performs the Spherical Pseudo-Mercator projection used by most web maps.
	Mercator = struct {
		ToWGS84 geo.Projection
	}{
		ToWGS84: func(p geo.Point) geo.Point {
			return geo.Point{
				180.0 * p[0] / earthRadiusPi,
				180.0 / math.Pi * (2*math.Atan(math.Exp(p[1]/geo.EarthRadius)) - math.Pi/2.0),
			}
		},
	}
	// WGS84 is what common uses lon/lat projection.
	WGS84 = struct {
		// ToMercator projections from WGS to Mercator, used by most web maps
		ToMercator geo.Projection
	}{
		ToMercator: func(g geo.Point) geo.Point {
			y := math.Log(math.Tan((90.0+g[1])*math.Pi/360.0)) * geo.EarthRadius
			return geo.Point{
				earthRadiusPi / 180.0 * g[0],
				math.Max(-earthRadiusPi, math.Min(y, earthRadiusPi)),
			}
		},
	}
)

// MercatorScaleFactor returns the mercator scaling factor for a given degree latitude.
func MercatorScaleFactor(g geo.Point) float64 {
	if g[1] < -90.0 || g[1] > 90.0 {
		panic(fmt.Sprintf("orb: latitude out of range, given %f", g[1]))
	}

	return 1.0 / math.Cos(g[1]/180.0*math.Pi)
}
