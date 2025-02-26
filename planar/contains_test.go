package planar

import "github.com/pchchv/geo"

func interpolate(a, b geo.Point, percent float64) geo.Point {
	return geo.Point{
		a[0] + percent*(b[0]-a[0]),
		a[1] + percent*(b[1]-a[1]),
	}
}
