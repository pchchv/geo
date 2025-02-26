package planar

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestRingContains(t *testing.T) {
	ring := geo.Ring{
		{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5},
		{2, 1}, {3, 1}, {3, 0}, {0, 0},
	}
	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+
	cases := []struct {
		name   string
		point  geo.Point
		result bool
	}{
		{
			name:   "in base",
			point:  geo.Point{1.5, 0.25},
			result: true,
		},
		{
			name:   "in right tower",
			point:  geo.Point{0.5, 0.75},
			result: true,
		},
		{
			name:   "in middle",
			point:  geo.Point{1.5, 0.75},
			result: false,
		},
		{
			name:   "in left tower",
			point:  geo.Point{2.5, 0.75},
			result: true,
		},
		{
			name:   "in tp middle",
			point:  geo.Point{1.5, 1.0},
			result: false,
		},
		{
			name:   "above",
			point:  geo.Point{2.5, 1.75},
			result: false,
		},
		{
			name:   "below",
			point:  geo.Point{2.5, -1.75},
			result: false,
		},
		{
			name:   "left",
			point:  geo.Point{-2.5, -0.75},
			result: false,
		},
		{
			name:   "right",
			point:  geo.Point{3.5, 0.75},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring.Reverse()
			if val := RingContains(ring, tc.point); val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}

			// should not care about orientation
			ring.Reverse()
			if val := RingContains(ring, tc.point); val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}
		})
	}

	// points should all be in
	for i, p := range ring {
		if !RingContains(ring, p) {
			t.Errorf("point index %d: should be inside", i)
		}
	}

	// on all the segments should be in.
	for i := 1; i < len(ring); i++ {
		c := interpolate(ring[i], ring[i-1], 0.5)
		if !RingContains(ring, c) {
			t.Errorf("index %d centroid: should be inside", i)
		}
	}

	// colinear with segments but outside
	for i := 1; i < len(ring); i++ {
		if p := interpolate(ring[i], ring[i-1], 5); RingContains(ring, p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}

		if p := interpolate(ring[i], ring[i-1], -5); RingContains(ring, p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}
	}
}

func TestPolygonContains(t *testing.T) {
	// should exclude holes
	p := geo.Polygon{
		{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
	}

	if !PolygonContains(p, geo.Point{1.5, 1.5}) {
		t.Errorf("should contain point")
	}

	// ring oriented same as outer ring
	p = append(p, geo.Ring{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}})
	if PolygonContains(p, geo.Point{1.5, 1.5}) {
		t.Errorf("should not contain point in hole")
	}

	p[1].Reverse() // oriented correctly as opposite of outer
	if PolygonContains(p, geo.Point{1.5, 1.5}) {
		t.Errorf("should not contain point in hole")
	}
}

func TestMultiPolygonContains(t *testing.T) {
	// should exclude holes
	mp := geo.MultiPolygon{
		{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
	}

	if !MultiPolygonContains(mp, geo.Point{0.5, 0.5}) {
		t.Errorf("should contain point")
	}

	if MultiPolygonContains(mp, geo.Point{1.5, 1.5}) {
		t.Errorf("should not contain point")
	}

	mp = append(mp, geo.Polygon{{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}}})
	if !MultiPolygonContains(mp, geo.Point{2.5, 0.5}) {
		t.Errorf("should contain point")
	}

	if MultiPolygonContains(mp, geo.Point{1.5, 0.5}) {
		t.Errorf("should not contain point")
	}
}

func interpolate(a, b geo.Point, percent float64) geo.Point {
	return geo.Point{
		a[0] + percent*(b[0]-a[0]),
		a[1] + percent*(b[1]-a[1]),
	}
}
