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

func TestMultiLineString(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{2, 2}}
	cases := []struct {
		name   string
		open   bool
		input  geo.MultiLineString
		output geo.MultiLineString
	}{
		{
			name: "regular closed bound clip",
			input: geo.MultiLineString{
				{{1, 1}, {2, 1}, {2, 2}, {3, 3}},
			},
			output: geo.MultiLineString{
				{{1, 1}, {2, 1}, {2, 2}, {2, 2}},
			},
		},
		{
			name: "open bound clip",
			open: true,
			input: geo.MultiLineString{
				{{1, 1}, {2, 1}, {2, 2}, {3, 3}},
			},
			output: geo.MultiLineString{
				{{1, 1}, {2, 1}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := MultiLineString(bound, tc.input, OpenBound(tc.open))
			if !result.Equal(tc.output) {
				t.Errorf("not equal")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestBound(t *testing.T) {
	cases := []struct {
		name string
		b1   geo.Bound
		b2   geo.Bound
		rs   geo.Bound
	}{
		{
			name: "normal intersection",
			b1:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			b2:   geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{4, 5}},
			rs:   geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{3, 4}},
		},
		{
			name: "1 contains 2",
			b1:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			b2:   geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{2, 3}},
			rs:   geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{2, 3}},
		},
		{
			name: "no overlap",
			b1:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			b2:   geo.Bound{Min: geo.Point{4, 5}, Max: geo.Point{5, 6}},
			rs:   geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{0, 0}}, // empty
		},
		{
			name: "same bound",
			b1:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			b2:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			rs:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
		},
		{
			name: "1 is empty",
			b1:   geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{0, 0}},
			b2:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
			rs:   geo.Bound{Min: geo.Point{0, 1}, Max: geo.Point{3, 4}},
		},
		{
			name: "both are empty",
			b1:   geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{0, 0}},
			b2:   geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{0, 0}},
			rs:   geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r1 := Bound(tc.b1, tc.b2)
			r2 := Bound(tc.b1, tc.b2)
			if tc.rs.IsEmpty() && (!r1.IsEmpty() || !r2.IsEmpty()) {
				t.Errorf("should be empty")
				t.Logf("%v", r1)
				t.Logf("%v", r2)
			}

			if !tc.rs.IsEmpty() {
				if !r1.Equal(tc.rs) {
					t.Errorf("r1 not equal")
					t.Logf("%v", r1)
					t.Logf("%v", tc.rs)
				}
				if !r2.Equal(tc.rs) {
					t.Errorf("r2 not equal")
					t.Logf("%v", r2)
					t.Logf("%v", tc.rs)
				}
			}
		})
	}
}

func TestGeometry(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{-1, -1}, Max: geo.Point{1, 1}}
	for _, g := range geo.AllGeometries {
		Geometry(bound, g)
	}

	cases := []struct {
		name   string
		input  geo.Geometry
		output geo.Geometry
	}{
		{
			name:   "only one multipoint in bound",
			input:  geo.MultiPoint{{0, 0}, {5, 5}},
			output: geo.Point{0, 0},
		},
		{
			name: "only one multilinestring in bound",
			input: geo.MultiLineString{
				{{0, 0}, {5, 5}},
				{{6, 6}, {7, 7}},
			},
			output: geo.LineString{{0, 0}, {1, 1}},
		},
		{
			name: "only one multipolygon in bound",
			input: geo.MultiPolygon{
				{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
				{{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}},
			},
			output: geo.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Geometry(bound, tc.input)
			if !geo.Equal(result, tc.output) {
				t.Errorf("not equal")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}
