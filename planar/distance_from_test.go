package planar

import (
	"testing"

	"github.com/pchchv/geo"
)

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
