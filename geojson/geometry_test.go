package geojson

import (
	"strings"
	"testing"

	"github.com/pchchv/geo"
)

func TestGeometry(t *testing.T) {
	for _, g := range geo.AllGeometries {
		NewGeometry(g)
	}
}

func TestGeometryMarshal(t *testing.T) {
	cases := []struct {
		name    string
		geom    geo.Geometry
		include string
	}{
		{
			name:    "point",
			geom:    geo.Point{},
			include: `"type":"Point"`,
		},
		{
			name:    "multi point",
			geom:    geo.MultiPoint{},
			include: `"type":"MultiPoint"`,
		},
		{
			name:    "linestring",
			geom:    geo.LineString{},
			include: `"type":"LineString"`,
		},
		{
			name:    "multi linestring",
			geom:    geo.MultiLineString{},
			include: `"type":"MultiLineString"`,
		},
		{
			name:    "polygon",
			geom:    geo.Polygon{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "multi polygon",
			geom:    geo.MultiPolygon{},
			include: `"type":"MultiPolygon"`,
		},
		{
			name:    "ring",
			geom:    geo.Ring{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "bound",
			geom:    geo.Bound{},
			include: `"type":"Polygon"`,
		},
		{
			name:    "collection",
			geom:    geo.Collection{geo.LineString{}},
			include: `"type":"GeometryCollection"`,
		},
		{
			name:    "collection2",
			geom:    geo.Collection{geo.Point{}, geo.Point{}},
			include: `"geometries":[`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NewGeometry(tc.geom).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}

			g := &Geometry{Coordinates: tc.geom}
			data, err = g.MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}
		})
	}
}
