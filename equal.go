package geo

import "fmt"

// Equal returns if the two geometrires are equal.
func Equal(g1, g2 Geometry) bool {
	if g1 == nil || g2 == nil {
		return g1 == g2
	}

	if g1.GeoJSONType() != g2.GeoJSONType() {
		return false
	}

	switch g1 := g1.(type) {
	case Point:
		return g1.Equal(g2.(Point))
	case MultiPoint:
		return g1.Equal(g2.(MultiPoint))
	case LineString:
		return g1.Equal(g2.(LineString))
	case MultiLineString:
		return g1.Equal(g2.(MultiLineString))
	case Ring:
		if g2, ok := g2.(Ring); !ok {
			return false
		} else {
			return g1.Equal(g2)
		}
	case Polygon:
		if g2, ok := g2.(Polygon); !ok {
			return false
		} else {
			return g1.Equal(g2)
		}
	case MultiPolygon:
		return g1.Equal(g2.(MultiPolygon))
	case Collection:
		return g1.Equal(g2.(Collection))
	case Bound:
		if g2, ok := g2.(Bound); !ok {
			return false
		} else {
			return g1.Equal(g2)
		}
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g1))
	}
}
