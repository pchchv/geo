package geojson

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
)

type geometry interface {
	Geometry() geo.Geometry
}

// This test makes sure the marshal-unmarshal loop does the same thing.
// The code and types here are complicated to avoid duplicate code.
func TestHelperTypes(t *testing.T) {
	cases := []struct {
		name   string
		geom   geo.Geometry
		helper interface{}
		output interface{}
	}{
		{
			name:   "point",
			geom:   geo.Point{1, 2},
			helper: Point(geo.Point{1, 2}),
			output: &Point{},
		},
		{
			name:   "multi point",
			geom:   geo.MultiPoint{{1, 2}, {3, 4}},
			helper: MultiPoint(geo.MultiPoint{{1, 2}, {3, 4}}),
			output: &MultiPoint{},
		},
		{
			name:   "linestring",
			geom:   geo.LineString{{1, 2}, {3, 4}},
			helper: LineString(geo.LineString{{1, 2}, {3, 4}}),
			output: &LineString{},
		},
		{
			name:   "multi linestring",
			geom:   geo.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			helper: MultiLineString(geo.MultiLineString{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}),
			output: &MultiLineString{},
		},
		{
			name:   "polygon",
			geom:   geo.Polygon{{{1, 2}, {3, 4}}},
			helper: Polygon(geo.Polygon{{{1, 2}, {3, 4}}}),
			output: &Polygon{},
		},
		{
			name:   "multi polygon",
			geom:   geo.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}}},
			helper: MultiPolygon(geo.MultiPolygon{{{{1, 2}, {3, 4}}}, {{{5, 6}, {7, 8}}}}),
			output: &MultiPolygon{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// check marshalling
			data, err := json.Marshal(tc.helper)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			geoData, err := json.Marshal(NewGeometry(tc.geom))
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !reflect.DeepEqual(data, geoData) {
				t.Errorf("should marshal the same")
				t.Log(string(data))
				t.Log(string(geoData))
			}

			// check unmarshalling
			if err = json.Unmarshal(data, tc.output); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			geom := &Geometry{}
			if err = json.Unmarshal(data, geom); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !geo.Equal(tc.output.(geometry).Geometry(), geom.Coordinates) {
				t.Errorf("should unmarshal the same")
				t.Log(tc.output)
				t.Log(geom.Coordinates)
			}

			// invalid json should return error
			if err = json.Unmarshal([]byte(`{invalid}`), tc.output); err == nil {
				t.Errorf("should return error for invalid json")
			}

			// not the correct type should return error.
			// non of they types directly supported are geometry collections.
			data, err = json.Marshal(NewGeometry(geo.Collection{geo.Point{}}))
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}

			if err = json.Unmarshal(data, tc.output); err == nil {
				t.Fatalf("should return error for invalid json")
			}
		})

		t.Run("bson "+tc.name, func(t *testing.T) {
			// check marshalling
			data, err := bson.Marshal(tc.helper)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if geoData, err := bson.Marshal(NewGeometry(tc.geom)); err != nil {
				t.Fatalf("marshal error: %v", err)
			} else if !reflect.DeepEqual(data, geoData) {
				t.Errorf("should marshal the same")
				t.Log(data)
				t.Log(geoData)
			}

			// check unmarshalling
			if err = bson.Unmarshal(data, tc.output); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			geom := &Geometry{}
			if err = bson.Unmarshal(data, geom); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !geo.Equal(tc.output.(geometry).Geometry(), geom.Coordinates) {
				t.Errorf("should unmarshal the same")
				t.Log(tc.output)
				t.Log(geom.Coordinates)
			}

			// invalid json should return error
			if err = bson.Unmarshal([]byte(`{invalid}`), tc.output); err == nil {
				t.Errorf("should return error for invalid bson")
			}

			// not the correct type should return error.
			// non of they types directly supported are geometry collections.
			data, err = bson.Marshal(NewGeometry(geo.Collection{geo.Point{}}))
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}

			if err = bson.Unmarshal(data, tc.output); err == nil {
				t.Fatalf("should return error for invalid bson")
			}
		})
	}
}
