package smartclip

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestNexts(t *testing.T) {
	for i, next := range nexts[geo.CW] {
		if next != -1 && i != nexts[geo.CCW][next] {
			t.Errorf("incorrect %d: %d != %d", i, i, nexts[geo.CCW][next])
		}
	}
}

func TestAroundBound(t *testing.T) {
	cases := []struct {
		name     string
		box      geo.Bound
		input    geo.Ring
		output   geo.Ring
		expected geo.Orientation
	}{
		{
			name:     "simple ccw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, -1}, {1, 1}},
			output:   geo.Ring{{-1, -1}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}},
			expected: geo.CCW,
		},
		{
			name:     "simple cw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, -1}, {1, 1}},
			output:   geo.Ring{{-1, -1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}},
			expected: geo.CW,
		},
		{
			name:     "wrap edge around whole box ccw",
			box:      geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{6, 6}},
			input:    geo.Ring{{1, 3}, {1, 2}},
			output:   geo.Ring{{1, 3}, {1, 2}, {1, 1}, {3.5, 1}, {6, 1}, {6, 3.5}, {6, 6}, {3.5, 6}, {1, 6}, {1, 3}},
			expected: geo.CCW,
		},
		{
			name:     "wrap around whole box ccw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}},
			output:   geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0.5}},
			expected: geo.CCW,
		},
		{
			name:     "wrap around whole box cw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, -0.5}, {0, -0.5}, {0, 0.5}, {-1, 0.5}},
			output:   geo.Ring{{-1, -0.5}, {0, -0.5}, {0, 0.5}, {-1, 0.5}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, -0.5}},
			expected: geo.CW,
		},
		{
			name:     "already cw with endpoints in same section",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}},
			output:   geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}, {-1, 0.5}},
			expected: geo.CW,
		},
		{
			name:     "cw but want ccw with endpoints in same section",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}},
			output:   geo.Ring{{-1, 0.5}, {0, 0.5}, {0, -0.5}, {-1, -0.5}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0.5}},
			expected: geo.CCW,
		},
		{
			name:     "one point on edge ccw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, 0.0}, {-0.5, -0.5}, {0, 0}, {-0.5, 0.5}, {-1, 0.0}},
			output:   geo.Ring{{-1, 0.0}, {-0.5, -0.5}, {0, 0}, {-0.5, 0.5}, {-1, 0.0}},
			expected: geo.CCW,
		},
		{
			name:     "one point on edge cw",
			box:      geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}},
			input:    geo.Ring{{-1, 0.0}, {-0.5, -0.5}, {0, 0}, {-0.5, 0.5}, {-1, 0.0}},
			output:   geo.Ring{{-1, 0.0}, {-0.5, -0.5}, {0, 0}, {-0.5, 0.5}, {-1, 0.0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}},
			expected: geo.CW,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := aroundBound(tc.box, tc.input, tc.expected)
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("does not match")
				t.Logf("%v", out)
				t.Logf("%v", tc.output)
			}
		})
	}
}
