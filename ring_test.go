package geo

import "testing"

func TestRing_Orientation(t *testing.T) {
	cases := []struct {
		name   string
		ring   Ring
		result Orientation
	}{
		{
			name:   "simple box, ccw",
			ring:   Ring{{0, 0}, {0.001, 0}, {0.001, 0.001}, {0, 0.001}, {0, 0}},
			result: CCW,
		},
		{
			name:   "simple box, cw",
			ring:   Ring{{0, 0}, {0, 0.001}, {0.001, 0.001}, {0.001, 0}, {0, 0}},
			result: CW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := tc.ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong orientation: %v != %v", val, tc.result)
			}

			// should work without redudant last point.
			ring := tc.ring[:len(tc.ring)-1]
			val = ring.Orientation()
			if val != tc.result {
				t.Errorf("wrong orientation: %v != %v", val, tc.result)
			}
		})
	}
}
