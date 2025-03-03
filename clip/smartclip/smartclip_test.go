package smartclip

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestSmartWrap(t *testing.T) {
	cases := []struct {
		name     string
		bound    geo.Bound
		rings    []geo.LineString
		expected geo.MultiPolygon
	}{
		{
			name:  "basic example",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{5, 5}},
			rings: []geo.LineString{
				{{0, 1}, {4, 1}, {4, 4}, {0, 4}},
				{{0, 3}, {3, 3}, {3, 2}, {0, 2}},
			},
			expected: geo.MultiPolygon{
				{{{0, 1}, {4, 1}, {4, 4}, {0, 4}, {0, 3}, {3, 3}, {3, 2}, {0, 2}, {0, 1}}},
			},
		},
		{
			name:  "two open one on each side of bound",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{5, 5}},
			rings: []geo.LineString{
				{{0, 2}, {2, 2}, {2, 3}, {0, 3}},
				{{5, 3}, {3, 3}, {3, 2}, {5, 2}},
			},
			expected: geo.MultiPolygon{
				{{{0, 2}, {2, 2}, {2, 3}, {0, 3}, {0, 2}}},
				{{{5, 3}, {3, 3}, {3, 2}, {5, 2}, {5, 3}}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := smartWrap(tc.bound, tc.rings, geo.CCW)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("incorrect ring")
				t.Logf("%v", result)
				t.Logf("%v", tc.expected)
			}
		})
	}
}

func TestSmartClip(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}}
	for _, g := range geo.AllGeometries {
		Geometry(bound, g, geo.CCW)
	}
}
