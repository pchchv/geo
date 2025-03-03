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

func TestPolygon(t *testing.T) {
	oneSix := geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{6, 6}}
	cases := []struct {
		name     string
		bound    geo.Bound
		input    geo.Polygon
		expected geo.MultiPolygon
	}{
		{
			name:  "with innner ring",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {5, 2}, {5, 5}, {0, 5}, {0, 2}},
				{{3, 3}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
			},
			expected: geo.MultiPolygon{{
				{{1, 2}, {5, 2}, {5, 5}, {1, 5}, {1, 2}},
				{{3, 3}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
			}},
		},
		{
			name:  "with innner ring that will share a side with the outer ring",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 5}, {0, 5}},
				{{0, 4}, {2, 4}, {2, 3}, {0, 3}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {3, 2}, {3, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 2}}},
			},
		},
		{
			name:  "with inner endpoints touching edge",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{9, 9}},
			input: geo.Polygon{
				{{0, 2}, {8, 2}, {8, 8}, {0, 8}},
				{{1, 4}, {2, 5}, {3, 4}, {2, 3}, {1, 4}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {8, 2}, {8, 8}, {1, 8}, {1, 4}, {2, 5}, {3, 4}, {2, 3}, {1, 4}, {1, 2}}},
			},
		},
		{
			name:  "with inner endpoints touching edge and closed inner",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{9, 9}},
			input: geo.Polygon{
				{{0, 2}, {8, 2}, {8, 8}, {0, 8}},
				{{6, 6}, {6, 7}, {7, 7}, {7, 6}, {6, 6}},
				{{1, 4}, {2, 5}, {3, 4}, {2, 3}, {1, 4}},
			},
			expected: geo.MultiPolygon{
				{
					{{1, 2}, {8, 2}, {8, 8}, {1, 8}, {1, 4}, {2, 5}, {3, 4}, {2, 3}, {1, 4}, {1, 2}},
					{{6, 6}, {6, 7}, {7, 7}, {7, 6}, {6, 6}},
				},
			},
		},
		{
			name:  "with inner interior point touching edge",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{9, 9}},
			input: geo.Polygon{
				{{0, 2}, {8, 2}, {8, 8}, {0, 8}},
				{{2, 5}, {3, 4}, {2, 3}, {1, 4}, {2, 5}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {8, 2}, {8, 8}, {1, 8}, {1, 4}, {2, 5}, {3, 4}, {2, 3}, {1, 4}, {1, 2}}},
			},
		},
		{
			name:  "with edge of inner along side",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 5}, {0, 5}},
				{{2, 4}, {2, 3}, {1, 3}, {1, 4}, {2, 4}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {3, 2}, {3, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 2}}},
			},
		},
		{
			name:  "with first edge of inner along side",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 5}, {0, 5}},
				{{1, 3}, {1, 4}, {2, 4}, {2, 3}, {1, 3}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {3, 2}, {3, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 2}}},
			},
		},
		{
			name:  "with last edge of inner along side",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 5}, {0, 5}},
				{{1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 4}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {3, 2}, {3, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 2}}},
			},
		},
		{
			name:  "with multiple edges of inner along side",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {4, 2}, {4, 7}},
				{{2, 5}, {1, 5}, {1, 6}, {2, 6}, {2, 5}},
			},
			expected: geo.MultiPolygon{
				{{{2, 6}, {2, 5}, {1, 5}, {1, 2}, {4, 2}, {4, 6}, {2, 6}}},
			},
		},
		{
			name:  "two non-continuous inner points on the boundary",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {4, 2}, {4, 7}},
				{{3, 5}, {2, 4}, {1, 5}, {2, 6}, {3, 5}},
			},
			expected: geo.MultiPolygon{
				{{{2, 6}, {3, 5}, {2, 4}, {1, 5}, {1, 2}, {4, 2}, {4, 6}, {2, 6}}},
				{{{1, 5}, {2, 6}, {1, 6}, {1, 5}}},
			},
		},
		{
			name:  "two inner points end up touching outer",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 7}},
				{{2, 4}, {1, 5}, {2, 6}, {3, 5}, {2, 4}},
			},
			expected: geo.MultiPolygon{
				{{{2, 6}, {3, 5}, {2, 4}, {1, 5}, {1, 2}, {3, 2}, {3, 6}, {2, 6}}},
				{{{1, 5}, {2, 6}, {1, 6}, {1, 5}}},
			},
		},
		{
			name:  "both inner and outer go outside and are open",
			bound: oneSix,
			input: geo.Polygon{
				{{0, 2}, {3, 2}, {3, 7}},
				{{1, 7}, {3, 5}, {2, 4}, {0, 6}},
			},
			expected: geo.MultiPolygon{
				{{{2, 6}, {3, 5}, {2, 4}, {1, 5}, {1, 2}, {3, 2}, {3, 6}, {2, 6}}},
			},
		},
	}

	for _, tc := range cases {
		for i, s := range flips {
			t.Run(tc.name+" - "+s, func(t *testing.T) {
				o := geo.CCW
				if i == 1 || i == 2 {
					o = geo.CW
				}
				bound := flipBound(i, tc.bound)

				input := tc.input.Clone()
				flipPolygon(i, input)

				expected := tc.expected.Clone()
				flipMultiPolygon(i, expected)

				result := Polygon(bound, input.Clone(), o)
				if !deepEqualMultiPolygon(result, expected) {
					t.Errorf("incorrect polygon")
					t.Logf("%v", result)
					t.Logf("%v", expected)
				}

				// should give same result if inputed as multi polygon
				result = MultiPolygon(bound, geo.MultiPolygon{input}, o)
				if !deepEqualMultiPolygon(result, expected) {
					t.Errorf("incorrect multipolygon")
					t.Logf("%v", result)
					t.Logf("%v", expected)
				}
			})
		}
	}
}

func TestClipMultiPolygon(t *testing.T) {
	oneSix := geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{6, 6}}
	cases := []struct {
		name     string
		bound    geo.Bound
		input    geo.MultiPolygon
		expected geo.MultiPolygon
	}{
		{
			name:  "two open rings",
			bound: oneSix,
			input: geo.MultiPolygon{
				{{{0, 2}, {2, 2}, {2, 7}, {8, 7}}},
				{{{8, 5}, {4, 5}, {4, 3}, {7, 3}}},
			},
			expected: geo.MultiPolygon{
				{{{6, 5}, {4, 5}, {4, 3}, {6, 3}, {6, 5}}},
				{{{1, 2}, {2, 2}, {2, 6}, {1, 6}, {1, 2}}},
			},
		},
		{
			name:  "two closed rings within the bound",
			bound: oneSix,
			input: geo.MultiPolygon{
				{{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}},
				{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}},
			},
			expected: geo.MultiPolygon{
				{{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}},
				{{{4, 4}, {5, 4}, {5, 5}, {4, 5}, {4, 4}}},
			},
		},
		{
			name:  "two open outer rings within each other",
			bound: oneSix,
			input: geo.MultiPolygon{
				{{{0, 2}, {3, 2}, {3, 5}, {0, 5}}},
				{{{0, 4}, {2, 4}, {2, 3}, {0, 3}}},
			},
			expected: geo.MultiPolygon{
				{{{1, 2}, {3, 2}, {3, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 2}}},
			},
		},
		{
			name:  "polygon with innner ring and single ring polygon",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{9, 9}},
			input: geo.MultiPolygon{
				{
					{{0, 2}, {5, 2}, {5, 5}, {0, 5}, {0, 2}},
					{{3, 3}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
				},
				{
					{{7, 7}, {8, 7}, {8, 8}, {7, 8}, {7, 7}},
				},
			},
			expected: geo.MultiPolygon{
				{
					{{1, 2}, {5, 2}, {5, 5}, {1, 5}, {1, 2}},
					{{3, 3}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
				},
				{
					{{7, 7}, {8, 7}, {8, 8}, {7, 8}, {7, 7}},
				},
			},
		},
	}

	for _, tc := range cases {
		for i, s := range flips {
			t.Run(tc.name+" - "+s, func(t *testing.T) {
				o := geo.CCW
				if i == 1 || i == 2 {
					o = geo.CW
				}
				bound := flipBound(i, tc.bound)

				input := tc.input.Clone()
				flipMultiPolygon(i, input)

				expected := tc.expected.Clone()
				flipMultiPolygon(i, expected)

				result := MultiPolygon(bound, input, o)
				if !deepEqualMultiPolygon(result, expected) {
					t.Errorf("incorrect multipolygon")
					t.Logf("%v", result)
					t.Logf("%v", expected)
				}
			})
		}
	}
}
