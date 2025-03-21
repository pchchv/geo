package geometries

import (
	"math"
	"testing"

	"github.com/pchchv/geo"
)

func TestSignedArea(t *testing.T) {
	area := 12392.029
	cases := []struct {
		name   string
		ring   geo.Ring
		result float64
	}{
		{
			name:   "simple box, ccw",
			ring:   geo.Ring{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: area,
		},
		{
			name:   "simple box, cc",
			ring:   geo.Ring{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: -area,
		},
		{
			name:   "even number of points",
			ring:   geo.Ring{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0.0004, 0.001}, {0, 0.001}, {0, 0}},
			result: area,
		},
		{
			name:   "3 points",
			ring:   geo.Ring{{0, 0}, {0.001, 0}, {0.001, 0.001}},
			result: area / 2.0,
		},
		{
			name:   "4 points",
			ring:   geo.Ring{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0}},
			result: area / 2.0,
		},
		{
			name:   "6 points",
			ring:   geo.Ring{{0.001, 0.001}, {0.002, 0.001}, {0.002, 0.0015}, {0.002, 0.002}, {0.001, 0.002}, {0.001, 0.001}},
			result: area,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := SignedArea(tc.ring)
			if math.Abs(val-tc.result) > 1 {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// should work without redudant last point
			if tc.ring[0] == tc.ring[len(tc.ring)-1] {
				tc.ring = tc.ring[:len(tc.ring)-1]
				val = SignedArea(tc.ring)
				if math.Abs(val-tc.result) > 1 {
					t.Errorf("wrong area: %v != %v", val, tc.result)
				}
			}
		})
	}
}

func TestArea(t *testing.T) {
	for _, g := range geo.AllGeometries {
		// should not panic with unsupported type
		Area(g)
	}
}
