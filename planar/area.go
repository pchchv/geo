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
