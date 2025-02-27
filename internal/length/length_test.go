package length

import (
	"math"
	"testing"

	"github.com/pchchv/geo"
)

func Distance(a, b geo.Point) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

func TestLength(t *testing.T) {
	for _, g := range geo.AllGeometries {
		// should not panic with unsupported type
		Length(g, Distance)
	}

	ls := geo.LineString{{0, 0}, {3, 0}, {3, 4}, {0, 0}}
	if l := Length(ls, Distance); l != 12 {
		t.Errorf("incorrect length: %v != %v", l, 12)
	}

	mls := geo.MultiLineString{
		{{0, 0}, {3, 0}, {3, 4}, {0, 0}},
		{{5, 0}, {5, 7}},
	}
	if l := Length(mls, Distance); l != 19 {
		t.Errorf("incorrect length: %v != %v", l, 19)
	}

	p := geo.Polygon{{{0, 0}, {3, 0}, {3, 4}, {0, 0}}}
	if l := Length(p, Distance); l != 12 {
		t.Errorf("incorrect length: %v != %v", l, 12)
	}

	mp := geo.MultiPolygon{
		{{{0, 0}, {3, 0}, {3, 4}, {0, 0}}},
		{{{5, 0}, {8, 0}, {8, 4}, {5, 0}}},
	}
	if l := Length(mp, Distance); l != 24 {
		t.Errorf("incorrect length: %v != %v", l, 24)
	}
}
