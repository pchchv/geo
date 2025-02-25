package planar

import "github.com/pchchv/geo"

func multiPointCentroid(mp geo.MultiPoint) geo.Point {
	if len(mp) == 0 {
		return geo.Point{}
	}

	x, y := 0.0, 0.0
	for _, p := range mp {
		x += p[0]
		y += p[1]
	}

	num := float64(len(mp))
	return geo.Point{x / num, y / num}
}
