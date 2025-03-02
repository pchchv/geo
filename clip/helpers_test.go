package clip

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  geo.Bound
		input  geo.Ring
		output geo.Ring
	}{
		{
			name:  "regular clip",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1.5, 1.5}},
			input: geo.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: geo.Ring{
				{1, 1}, {1.5, 1}, {1.5, 1.5}, {1, 1.5}, {1, 1},
			},
		},
		{
			name:  "bound to the top",
			bound: geo.Bound{Min: geo.Point{-1, 3}, Max: geo.Point{3, 4}},
			input: geo.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: geo.Ring{},
		},
		{
			name:  "bound in lower left",
			bound: geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{0, 0}},
			input: geo.Ring{
				{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
			},
			output: geo.Ring{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Ring(tc.bound, tc.input)
			if !result.Equal(tc.output) {
				t.Errorf("not equal")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}
