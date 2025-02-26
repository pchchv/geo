package planar

import (
	"math"
	"testing"

	"github.com/pchchv/geo"
)

const epsilon = 1e-6

func TestDistanceFromSegment(t *testing.T) {
	a := geo.Point{0, 0}
	b := geo.Point{0, 10}

	cases := []struct {
		name   string
		point  geo.Point
		result float64
	}{
		{
			name:   "point in middle",
			point:  geo.Point{1, 5},
			result: 1,
		},
		{
			name:   "on line",
			point:  geo.Point{0, 2},
			result: 0,
		},
		{
			name:   "past start",
			point:  geo.Point{0, -5},
			result: 5,
		},
		{
			name:   "past end",
			point:  geo.Point{0, 13},
			result: 3,
		},
		{
			name:   "triangle",
			point:  geo.Point{3, 4},
			result: 3,
		},
		{
			name:   "triangle off end",
			point:  geo.Point{3, -4},
			result: 5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if v := DistanceFromSegment(a, b, tc.point); v != tc.result {
				t.Errorf("incorrect distance: %v != %v", v, tc.result)
			}
		})
	}
}

func TestDistanceFromWithIndex(t *testing.T) {
	for _, g := range geo.AllGeometries {
		DistanceFromWithIndex(g, geo.Point{})
	}
}

func TestDistanceFrom_LineString(t *testing.T) {
	ls := geo.LineString{{0, 0}, {0, 3}, {4, 3}, {4, 0}}
	cases := []struct {
		name   string
		point  geo.Point
		result float64
	}{
		{
			point:  geo.Point{4.5, 1.5},
			result: 0.5,
		},
		{
			point:  geo.Point{0.4, 1.5},
			result: 0.4,
		},
		{
			point:  geo.Point{-0.3, 1.5},
			result: 0.3,
		},
		{
			point:  geo.Point{0.3, 2.8},
			result: 0.2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := DistanceFrom(ls, tc.point)
			if math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestDistanceFrom_Polygon(t *testing.T) {
	r1 := geo.Ring{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}}
	r2 := geo.Ring{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}}
	poly := geo.Polygon{r1, r2}
	cases := []struct {
		name   string
		point  geo.Point
		result float64
	}{
		{
			name:   "outside",
			point:  geo.Point{-1, 2},
			result: 1,
		},
		{
			name:   "inside",
			point:  geo.Point{0.4, 2},
			result: 0.4,
		},
		{
			name:   "in hole",
			point:  geo.Point{1.3, 1.4},
			result: 0.3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if d := DistanceFrom(poly, tc.point); math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestDistanceFrom_MultiPoint(t *testing.T) {
	mp := geo.MultiPoint{{0.0}, {1, 1}, {2, 2}}
	fromPoint := geo.Point{3, 2}
	if distance := DistanceFrom(mp, fromPoint); distance != 1 {
		t.Errorf("distance incorrect: %v != %v", distance, 1)
	}
}
