package geojson

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
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
				t.Fatalf("marshal error: %e", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}

			g := &Geometry{Coordinates: tc.geom}
			data, err = g.MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %e", err)
			}

			if !strings.Contains(string(data), tc.include) {
				t.Errorf("does not contain substring")
				t.Log(string(data))
			}
		})
	}
}

func TestGeometryUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		geom geo.Geometry
	}{
		{
			name: "point",
			geom: geo.Point{1, 2},
		},
		{
			name: "multi point",
			geom: geo.MultiPoint{{1, 2}, {3, 4}},
		},
		{
			name: "linestring",
			geom: geo.LineString{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name: "multi linestring",
			geom: geo.MultiLineString{},
		},
		{
			name: "polygon",
			geom: geo.Polygon{},
		},
		{
			name: "multi polygon",
			geom: geo.MultiPolygon{},
		},
		{
			name: "collection",
			geom: geo.Collection{geo.LineString{{1, 2}, {3, 4}}, geo.Point{5, 6}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NewGeometry(tc.geom).MarshalJSON()
			if err != nil {
				t.Fatalf("marshal error: %e", err)
			}

			// unmarshal
			g, err := UnmarshalGeometry(data)
			if err != nil {
				t.Errorf("unmarshal error: %e", err)
			}

			if g.Type != tc.geom.GeoJSONType() {
				t.Errorf("incorrenct type: %v != %v", g.Type, tc.geom.GeoJSONType())
			}

			if !geo.Equal(g.Geometry(), tc.geom) {
				t.Errorf("incorrect geometry")
				t.Logf("%[1]T, %[1]v", g.Geometry())
				t.Log(tc.geom)
			}
		})
	}

	// invalid type
	if _, err := UnmarshalGeometry([]byte(`{
		"type": "arc",
		"coordinates": [[0, 0]]
	}`)); err == nil {
		t.Errorf("should return error for invalid type")
	} else if !strings.Contains(err.Error(), "invalid geometry") {
		t.Errorf("incorrect error: %e", err)
	}

	// invalid json
	// truncated
	if _, err := UnmarshalGeometry([]byte(`{"type": "arc",`)); err == nil {
		t.Errorf("should return error for invalid json")
	}

	g := &Geometry{}
	// truncated
	if err := g.UnmarshalJSON([]byte(`{"type": "arc",`)); err == nil {
		t.Errorf("should return error for invalid json")
	}

	// invalid type (null)
	if _, err := UnmarshalGeometry([]byte(`null`)); err == nil {
		t.Errorf("should return error for invalid type")
	} else if !strings.Contains(err.Error(), "invalid geometry") {
		t.Errorf("incorrect error: %e", err)
	}
}

func TestGeometryUnmarshal_errors(t *testing.T) {
	cases := []struct {
		name string
		data string
	}{
		{
			name: "point",
			data: `{"type":"Point","coordinates":1}`,
		},
		{
			name: "multi point",
			data: `{"type":"MultiPoint","coordinates":2}`,
		},
		{
			name: "linestring",
			data: `{"type":"LineString","coordinates":3}`,
		},
		{
			name: "multi linestring",
			data: `{"type":"MultiLineString","coordinates":4}`,
		},
		{
			name: "polygon",
			data: `{"type":"Polygon","coordinates":10.2}`,
		},
		{
			name: "multi polygon",
			data: `{"type":"MultiPolygon","coordinates":{}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalGeometry([]byte(tc.data)); err == nil {
				t.Errorf("expected error, got nothing")
			}
		})
	}
}

func TestGeometryMarshalJSON_null(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		type S struct {
			GeoJSON *Geometry `json:"geojson"`
		}

		var s S
		if err := json.Unmarshal([]byte(`{"geojson": null}`), &s); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		} else if s.GeoJSON != nil {
			t.Errorf("should be nil, got: %v", s)
		}
	})

	t.Run("feature with null geometry", func(t *testing.T) {
		type S struct {
			GeoJSON *Feature `json:"geojson"`
		}

		var s S
		if err := json.Unmarshal([]byte(`{"geojson": {"type":"Feature","geometry":null,"properties":null}}`), &s); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		} else if s.GeoJSON.Geometry != nil {
			t.Errorf("should be nil, got: %v", s)
		}
	})
}

func BenchmarkGeometryMarshalJSON(b *testing.B) {
	ls := geo.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, geo.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(g); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkGeometryUnmarshalJSON(b *testing.B) {
	ls := geo.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, geo.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}
	data, err := json.Marshal(g)
	if err != nil {
		b.Fatalf("marshal error: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, g); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkGeometryMarshalBSON(b *testing.B) {
	ls := geo.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, geo.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bson.Marshal(g); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkGeometryUnmarshalBSON(b *testing.B) {
	ls := geo.LineString{}
	for i := 0.0; i < 1000; i++ {
		ls = append(ls, geo.Point{i * 3.45, i * -58.4})
	}

	g := &Geometry{Coordinates: ls}
	data, err := bson.Marshal(g)
	if err != nil {
		b.Fatalf("marshal error: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := bson.Unmarshal(data, g); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}
