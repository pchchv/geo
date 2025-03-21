package simplifier

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestDoubleTriangleArea(t *testing.T) {
	expected := 30.0
	ls := geo.LineString{{2, 5}, {5, 1}, {-4, 3}}
	cases := []struct {
		name       string
		i1, i2, i3 int
	}{
		{
			name: "order 1",
			i1:   0, i2: 1, i3: 2,
		},
		{
			name: "order 2",
			i1:   0, i2: 2, i3: 1,
		},
		{
			name: "order 3",
			i1:   1, i2: 2, i3: 0,
		},
		{
			name: "order 4",
			i1:   1, i2: 0, i3: 2,
		},
		{
			name: "order 5",
			i1:   2, i2: 0, i3: 1,
		},
		{
			name: "order 6",
			i1:   2, i2: 1, i3: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			area := doubleTriangleArea(ls, tc.i1, tc.i2, tc.i3)
			if area != expected {
				t.Errorf("incorrect area: %v != %v", area, expected)
			}
		})
	}
}

func TestVisvalingamThreshold(t *testing.T) {
	cases := []struct {
		name      string
		threshold float64
		ls        geo.LineString
		expected  geo.LineString
		indexMap  []int
	}{
		{
			name:      "no reduction",
			threshold: 0.9,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected:  geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			indexMap:  []int{0, 1, 2, 3, 4},
		},
		{
			name:      "reduction",
			threshold: 1.1,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected:  geo.LineString{{0, 0}, {0, 4}},
			indexMap:  []int{0, 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, im := VisvalingamThreshold(tc.threshold).simplify(tc.ls, false, true)
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

func TestVisvalingamThreshold_area(t *testing.T) {
	cases := []struct {
		name     string
		r        geo.Ring
		expected geo.Ring
		indexMap []int
	}{
		{
			name:     "reduction",
			r:        geo.Ring{{0, 0}, {1, 0}, {1, 0.5}, {1, 1}, {0.5, 1}, {0, 1}},
			expected: geo.Ring{{0, 0}, {1, 0}, {0, 1}},
			indexMap: []int{0, 1, 5},
		},
		{
			name:     "reduction closed",
			r:        geo.Ring{{0, 0}, {1, 0}, {1, 0.5}, {1, 1}, {0.5, 1}, {0, 1}, {0, 0}},
			expected: geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			indexMap: []int{0, 1, 3, 6},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, im := VisvalingamThreshold(10).simplify(geo.LineString(tc.r), true, true)
			if !v.Equal(geo.LineString(tc.expected)) {
				t.Log(v)
				t.Log(tc.expected)
				t.Errorf("incorrect ring")
			}

			if !reflect.DeepEqual(im, tc.indexMap) {
				t.Log(im)
				t.Log(tc.indexMap)
				t.Errorf("incorrect index map")
			}
		})
	}
}

func TestVisvalingamKeep(t *testing.T) {
	cases := []struct {
		name     string
		keep     int
		ls       geo.LineString
		expected geo.LineString
		indexMap []int
	}{
		{
			name:     "keep 6",
			keep:     6,
			ls:       geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected: geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			indexMap: []int{0, 1, 2, 3, 4},
		},
		{
			name:     "keep 5",
			keep:     5,
			ls:       geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected: geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			indexMap: []int{0, 1, 2, 3, 4},
		},
		{
			name:     "keep 4",
			keep:     4,
			ls:       geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected: geo.LineString{{0, 0}, {0, 2}, {1, 3}, {0, 4}},
			indexMap: []int{0, 2, 3, 4},
		},
		{
			name:     "keep 3",
			keep:     3,
			ls:       geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected: geo.LineString{{0, 0}, {0, 2}, {0, 4}},
			indexMap: []int{0, 2, 4},
		},
		{
			name:     "keep 2",
			keep:     2,
			ls:       geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected: geo.LineString{{0, 0}, {0, 4}},
			indexMap: []int{0, 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, im := VisvalingamKeep(tc.keep).simplify(tc.ls, false, true)
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

func TestVisvalingam(t *testing.T) {
	cases := []struct {
		name      string
		threshold float64
		keep      int
		ls        geo.LineString
		expected  geo.LineString
		indexMap  []int
	}{
		{
			name:      "keep nothing",
			threshold: 1.1,
			keep:      0,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}},
			expected:  geo.LineString{{0, 0}, {0, 2}},
			indexMap:  []int{0, 2},
		},
		{
			name:      "keep everything",
			threshold: 1.1,
			keep:      3,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}},
			expected:  geo.LineString{{0, 0}, {1, 1}, {0, 2}},
			indexMap:  []int{0, 1, 2},
		},
		{
			name:      "not meeting threshold",
			threshold: 0.9,
			keep:      0,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}},
			expected:  geo.LineString{{0, 0}, {1, 1}, {0, 2}},
			indexMap:  []int{0, 1, 2},
		},
		{
			name:      "5 points keep nothing",
			threshold: 1.1,
			keep:      0,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected:  geo.LineString{{0, 0}, {0, 4}},
			indexMap:  []int{0, 4},
		},
		{
			name:      "5 points keep everything",
			threshold: 1.1,
			keep:      5,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected:  geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			indexMap:  []int{0, 1, 2, 3, 4},
		},
		{
			name:      "5 points reduce to limit",
			threshold: 1.1,
			keep:      3,
			ls:        geo.LineString{{0, 0}, {1, 1}, {0, 2}, {1, 3}, {0, 4}},
			expected:  geo.LineString{{0, 0}, {0, 2}, {0, 4}},
			indexMap:  []int{0, 2, 4},
		},
		{
			name:      "removes colinear points",
			threshold: 0.1,
			keep:      0,
			ls:        geo.LineString{{0, 0}, {0, 1}, {0, 2}},
			expected:  geo.LineString{{0, 0}, {0, 2}},
			indexMap:  []int{0, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, im := Visvalingam(tc.threshold, tc.keep).simplify(tc.ls, false, true)
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
