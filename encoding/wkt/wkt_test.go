package wkt

import (
	"bytes"
	"testing"

	"github.com/pchchv/geo"
)

func TestMarshal(t *testing.T) {
	cases := []struct {
		name     string
		geo      geo.Geometry
		expected []byte
	}{
		{
			name:     "point",
			geo:      geo.Point{1, 2},
			expected: []byte("POINT(1 2)"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := Marshal(tc.geo)
			if !bytes.Equal(v, tc.expected) {
				t.Log(string(v))
				t.Log(string(tc.expected))
				t.Errorf("incorrect wkt marshalling")
			}
		})
	}
}

func TestMarshalString(t *testing.T) {
	cases := []struct {
		name     string
		geo      geo.Geometry
		expected string
	}{
		{
			name:     "point",
			geo:      geo.Point{1, 2},
			expected: "POINT(1 2)",
		},
		{
			name:     "multipoint",
			geo:      geo.MultiPoint{{1, 2}, {0.5, 1.5}},
			expected: "MULTIPOINT((1 2),(0.5 1.5))",
		},
		{
			name:     "multipoint empty",
			geo:      geo.MultiPoint{},
			expected: "MULTIPOINT EMPTY",
		},
		{
			name:     "linestring",
			geo:      geo.LineString{{1, 2}, {0.5, 1.5}},
			expected: "LINESTRING(1 2,0.5 1.5)",
		},
		{
			name:     "linestring empty",
			geo:      geo.LineString{},
			expected: "LINESTRING EMPTY",
		},
		{
			name:     "multilinestring",
			geo:      geo.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			expected: "MULTILINESTRING((1 2,3 4),(5 6,7 8))",
		},
		{
			name:     "multilinestring empty",
			geo:      geo.MultiLineString{},
			expected: "MULTILINESTRING EMPTY",
		},
		{
			name:     "ring",
			geo:      geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			expected: "POLYGON((0 0,1 0,1 1,0 0))",
		},
		{
			name:     "polygon",
			geo:      geo.Polygon{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			expected: "POLYGON((1 2,3 4),(5 6,7 8))",
		},
		{
			name:     "polygon empty",
			geo:      geo.Polygon{},
			expected: "POLYGON EMPTY",
		},
		{
			name:     "multipolygon",
			geo:      geo.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}, {{1, 2}, {5, 4}}}},
			expected: "MULTIPOLYGON(((1 2,3 4)),((5 6,7 8),(1 2,5 4)))",
		},
		{
			name:     "multipolygon empty",
			geo:      geo.MultiPolygon{},
			expected: "MULTIPOLYGON EMPTY",
		},
		{
			name:     "collection",
			geo:      geo.Collection{geo.Point{1, 2}, geo.LineString{{3, 4}, {5, 6}}},
			expected: "GEOMETRYCOLLECTION(POINT(1 2),LINESTRING(3 4,5 6))",
		},
		{
			name:     "collection empty",
			geo:      geo.Collection{},
			expected: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name:     "bound",
			geo:      geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 2}},
			expected: "POLYGON((0 0,1 0,1 2,0 2,0 0))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := MarshalString(tc.geo)
			if v != tc.expected {
				t.Log(v)
				t.Log(tc.expected)
				t.Errorf("incorrect wkt marshalling")
			}
		})
	}
}
