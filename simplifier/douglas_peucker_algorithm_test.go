package simplifier

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestDouglasPeucker(t *testing.T) {
	cases := []struct {
		name      string
		threshold float64
		ls        geo.LineString
		expected  geo.LineString
		indexMap  []int
	}{
		{
			name:      "no reduction",
			threshold: 0.1,
			ls:        geo.LineString{{0, 0}, {0.5, 0.2}, {1, 0}},
			expected:  geo.LineString{{0, 0}, {0.5, 0.2}, {1, 0}},
			indexMap:  []int{0, 1, 2},
		},
		{
			name:      "reduction",
			threshold: 1.1,
			ls:        geo.LineString{{0, 0}, {0.5, 0.2}, {1, 0}},
			expected:  geo.LineString{{0, 0}, {1, 0}},
			indexMap:  []int{0, 2},
		},
		{
			name:      "removes coplanar points",
			threshold: 0,
			ls:        geo.LineString{{0, 0}, {0, 1}, {0, 2}},
			expected:  geo.LineString{{0, 0}, {0, 2}},
			indexMap:  []int{0, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, im := DouglasPeucker(tc.threshold).simplify(tc.ls, false, true)
			if !v.Equal(tc.expected) {
				t.Log(v)
				t.Log(tc.expected)
				t.Errorf("incorrect line")
			}

			if !reflect.DeepEqual(im, tc.indexMap) {
				t.Log(im)
				t.Log(tc.indexMap)
				t.Errorf("incorrect index map")
			}
		})
	}
}
