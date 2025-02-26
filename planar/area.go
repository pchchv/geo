package planar

import (
	"math"

	"github.com/pchchv/geo"
)

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

func ringCentroidArea(r geo.Ring) (centroid geo.Point, area float64) {
	if len(r) == 0 {
		return geo.Point{}, 0
	}

	// implicitly move everything to near the origin to help with roundoff
	offsetX := r[0][0]
	offsetY := r[0][1]
	for i := 1; i < len(r)-1; i++ {
		a := (r[i][0]-offsetX)*(r[i+1][1]-offsetY) - (r[i+1][0]-offsetX)*(r[i][1]-offsetY)
		area += a
		centroid[0] += (r[i][0] + r[i+1][0] - 2*offsetX) * a
		centroid[1] += (r[i][1] + r[i+1][1] - 2*offsetY) * a
	}

	if area == 0 {
		return r[0], 0
	}

	// no need to deal with first and last vertex since we
	// "moved" that point the origin (multiply by 0 == 0)
	area /= 2
	centroid[0] /= 6 * area
	centroid[1] /= 6 * area
	centroid[0] += offsetX
	centroid[1] += offsetY

	return centroid, area
}

func lineStringCentroidDist(ls geo.LineString) (geo.Point, float64) {
	if len(ls) == 0 {
		return geo.Point{}, math.Inf(1)
	}

	var dist float64
	offset := ls[0] // implicitly move everything to near the origin to help with roundoff
	point := geo.Point{}
	for i := 0; i < len(ls)-1; i++ {
		p1 := geo.Point{
			ls[i][0] - offset[0],
			ls[i][1] - offset[1],
		}

		p2 := geo.Point{
			ls[i+1][0] - offset[0],
			ls[i+1][1] - offset[1],
		}

		d := Distance(p1, p2)
		point[0] += (p1[0] + p2[0]) / 2.0 * d
		point[1] += (p1[1] + p2[1]) / 2.0 * d
		dist += d
	}

	if dist == 0 {
		return ls[0], 0
	}

	point[0] /= dist
	point[1] /= dist
	point[0] += ls[0][0]
	point[1] += ls[0][1]
	return point, dist
}
