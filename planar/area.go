package planar

import (
	"fmt"
	"math"

	"github.com/pchchv/geo"
)

// CentroidArea returns the centroid and the area in the 2d plane.
// Since area is needed for the centroid, returns both values.
// Polygon area will always be >= zero.
// Ring area can be negative if it is wound clockwise.
func CentroidArea(g geo.Geometry) (geo.Point, float64) {
	switch g := g.(type) {
	case nil:
		return geo.Point{}, 0
	case geo.Point:
		return multiPointCentroid(geo.MultiPoint{g}), 0
	case geo.MultiPoint:
		return multiPointCentroid(g), 0
	case geo.LineString:
		return multiLineStringCentroid(geo.MultiLineString{g}), 0
	case geo.MultiLineString:
		return multiLineStringCentroid(g), 0
	case geo.Ring:
		return ringCentroidArea(g)
	case geo.Polygon:
		return polygonCentroidArea(g)
	case geo.MultiPolygon:
		return multiPolygonCentroidArea(g)
	case geo.Collection:
		return collectionCentroidArea(g)
	case geo.Bound:
		return CentroidArea(g.ToRing())
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}

// Area returns the area of the geometry in the 2d plane.
// Needs to be finalized:
// want to implement area calculation without centroid calculation.
func Area(g geo.Geometry) float64 {
	_, a := CentroidArea(g)
	return a
}

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

func multiLineStringCentroid(mls geo.MultiLineString) (point geo.Point) {
	if len(mls) == 0 {
		return geo.Point{}
	}

	var dist float64
	var validCount int
	for _, ls := range mls {
		c, d := lineStringCentroidDist(ls)
		if d == math.Inf(1) {
			continue
		}

		dist += d
		validCount++
		if d == 0 {
			d = 1.0
		}

		point[0] += c[0] * d
		point[1] += c[1] * d
	}

	if validCount == 0 {
		return geo.Point{}
	}

	if dist == math.Inf(1) || dist == 0.0 {
		point[0] /= float64(validCount)
		point[1] /= float64(validCount)
		return
	}

	point[0] /= dist
	point[1] /= dist
	return
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

func lineStringCentroidDist(ls geo.LineString) (point geo.Point, dist float64) {
	if len(ls) == 0 {
		return geo.Point{}, math.Inf(1)
	}

	// implicitly move everything to near the origin to help with roundoff
	offset := ls[0]
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
	return
}

func polygonCentroidArea(p geo.Polygon) (geo.Point, float64) {
	if len(p) == 0 {
		return geo.Point{}, 0
	}

	centroid, area := ringCentroidArea(p[0])
	area = math.Abs(area)
	if len(p) == 1 {
		if area == 0 {
			c, _ := lineStringCentroidDist(geo.LineString(p[0]))
			return c, 0
		}
		return centroid, area
	}

	var holeArea float64
	weightedHoleCentroid := geo.Point{}
	for i := 1; i < len(p); i++ {
		hc, ha := ringCentroidArea(p[i])
		ha = math.Abs(ha)

		holeArea += ha
		weightedHoleCentroid[0] += hc[0] * ha
		weightedHoleCentroid[1] += hc[1] * ha
	}

	totalArea := area - holeArea
	if totalArea == 0 {
		c, _ := lineStringCentroidDist(geo.LineString(p[0]))
		return c, 0
	}

	centroid[0] = (area*centroid[0] - weightedHoleCentroid[0]) / totalArea
	centroid[1] = (area*centroid[1] - weightedHoleCentroid[1]) / totalArea
	return centroid, totalArea
}

func multiPolygonCentroidArea(mp geo.MultiPolygon) (point geo.Point, area float64) {
	for _, p := range mp {
		c, a := polygonCentroidArea(p)
		point[0] += c[0] * a
		point[1] += c[1] * a
		area += a
	}

	if area == 0 {
		return geo.Point{}, 0
	}

	point[0] /= area
	point[1] /= area
	return
}

func collectionCentroidArea(c geo.Collection) (point geo.Point, area float64) {
	max := maxDim(c)
	for _, g := range c {
		if g.Dimensions() != max {
			continue
		}

		c, a := CentroidArea(g)
		point[0] += c[0] * a
		point[1] += c[1] * a
		area += a
	}

	if area == 0 {
		return geo.Point{}, 0
	}

	point[0] /= area
	point[1] /= area
	return
}

func maxDim(c geo.Collection) (max int) {
	for _, g := range c {
		if d := g.Dimensions(); d > max {
			max = d
		}
	}

	return
}
