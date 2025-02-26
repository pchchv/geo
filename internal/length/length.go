package length

import (
	"fmt"

	"github.com/pchchv/geo"
)

// Length returns the length of the boundary of the geometry using 2d euclidean geometry.
func Length(g geo.Geometry, df geo.DistanceFunc) (sum float64) {
	switch g := g.(type) {
	case nil:
		return 0
	case geo.Point:
		return 0
	case geo.MultiPoint:
		return 0
	case geo.LineString:
		return lineStringLength(g, df)
	case geo.MultiLineString:
		for _, ls := range g {
			sum += lineStringLength(ls, df)
		}
		return
	case geo.Ring:
		return lineStringLength(geo.LineString(g), df)
	case geo.Polygon:
		return polygonLength(g, df)
	case geo.MultiPolygon:
		for _, p := range g {
			sum += polygonLength(p, df)
		}
		return
	case geo.Collection:
		for _, c := range g {
			sum += Length(c, df)
		}
		return
	case geo.Bound:
		return Length(g.ToRing(), df)
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}

func lineStringLength(ls geo.LineString, df geo.DistanceFunc) (sum float64) {
	for i := 1; i < len(ls); i++ {
		sum += df(ls[i], ls[i-1])
	}

	return
}

func polygonLength(p geo.Polygon, df geo.DistanceFunc) (sum float64) {
	for _, r := range p {
		sum += lineStringLength(geo.LineString(r), df)
	}

	return
}
