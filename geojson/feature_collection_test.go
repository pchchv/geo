package geojson

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewFeatureCollection(t *testing.T) {
	fc := NewFeatureCollection()
	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}
}

func TestFeatureCollectionMarshal(t *testing.T) {
	fc := NewFeatureCollection()
	fc.Features = nil
	if blob, err := json.Marshal(fc); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}

func TestFeatureCollectionMarshalJSON(t *testing.T) {
	fc := NewFeatureCollection()
	if blob, err := fc.MarshalJSON(); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}

func TestFeatureCollectionMarshalJSON_null(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		type S struct {
			GeoJSON *FeatureCollection `json:"geojson"`
		}

		var s S
		if err := json.Unmarshal([]byte(`{"geojson": null}`), &s); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		}

		if s.GeoJSON != nil {
			t.Errorf("should be nil, got: %v", s)
		}
	})

	t.Run("non-pointer", func(t *testing.T) {
		type S struct {
			GeoJSON FeatureCollection `json:"geojson"`
		}

		var s S
		if err := json.Unmarshal([]byte(`{"geojson": null}`), &s); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		}

		if !reflect.DeepEqual(s.GeoJSON, FeatureCollection{}) {
			t.Errorf("should be empty, got: %v", s)
		}
	})
}

func TestFeatureCollectionMarshal_BBox(t *testing.T) {
	fc := NewFeatureCollection()
	fc.BBox = nil
	if blob, err := json.Marshal(fc); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if bytes.Contains(blob, []byte(`"bbox"`)) {
		t.Errorf("should not contain bbox attribute if empty")
	}

	// with a bbox
	fc.BBox = []float64{1, 2, 3, 4}
	if blob, err := json.Marshal(fc); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"bbox":[1,2,3,4]`)) {
		t.Errorf("did not marshal bbox correctly: %v", string(blob))
	}
}

func TestFeatureCollectionMarshalValue(t *testing.T) {
	fc := NewFeatureCollection()
	fc.Features = nil
	if blob, err := json.Marshal(*fc); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"features":[]`)) {
		t.Errorf("json should set features object to at least empty array")
	}
}

func TestFeatureCollectionMarshalJSON_extraMembers(t *testing.T) {
	rawJSON := `
	  { "type": "FeatureCollection",
		"foo": "bar",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	     ]
	  }`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature collection without issue, err %e", err)
	}

	if v := fc.ExtraMembers.MustString("foo", ""); v != "bar" {
		t.Errorf("missing extra: foo: %v", v)
	}

	if data, err := fc.MarshalJSON(); err != nil {
		t.Fatalf("unable to marshal: %e", err)
	} else if !bytes.Contains(data, []byte(`"foo":"bar"`)) {
		t.Fatalf("extras not in marshalled data")
	}
}

func TestFeatureCollection_MarshalBSON(t *testing.T) {
	cases := []struct {
		name string
		geo  geo.Geometry
	}{
		{
			name: "point",
			geo:  geo.Point{1, 2},
		},
		{
			name: "multi point",
			geo:  geo.MultiPoint{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name: "line string",
			geo:  geo.LineString{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name: "multi line string",
			geo:  geo.MultiLineString{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}}},
		},
		{
			name: "polygon",
			geo:  geo.Polygon{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}}},
		},
		{
			name: "multi polygon",
			geo: geo.MultiPolygon{
				{
					{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}},
					{{9, 8}, {7, 6}, {5, 4}}, {{3, 2}, {1, 0}},
				},
				{
					{{9, 8}, {7, 6}, {5, 4}}, {{3, 2}, {1, 0}},
					{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}},
				},
			},
		},
		{
			name: "geometry collection",
			geo: geo.Collection{
				geo.Point{1, 2},
				geo.MultiPoint{{1, 2}, {3, 4}, {5, 6}},
				geo.LineString{{1, 2}, {3, 4}, {5, 6}},
				geo.MultiLineString{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}}},
				geo.Polygon{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 8}, {7, 6}}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fc := NewFeatureCollection()
			fc.Append(NewFeature(tc.geo))

			data, err := bson.Marshal(fc)
			if err != nil {
				t.Fatalf("unable to marshal feature collection: %e", err)
			}

			nfc := NewFeatureCollection()
			if err = bson.Unmarshal(data, &nfc); err != nil {
				t.Fatalf("unable to unmarshal feature collection: %e", err)
			}

			if nfc.Type != "FeatureCollection" {
				t.Errorf("feature collection type not set: %v", nfc.Type)
			}

			if nfc.Features[0].Geometry.GeoJSONType() != tc.geo.GeoJSONType() {
				t.Errorf("incorrect geometry type: %v != %v", nfc.Features[0].Geometry.GeoJSONType(), tc.geo.GeoJSONType())
			}

			if !geo.Equal(nfc.Features[0].Geometry, tc.geo) {
				t.Errorf("incorrect geometry: %v != %v", nfc.Features[0].Geometry, tc.geo)
			}

			if nfc.Features[0].Type != "Feature" {
				t.Errorf("feature type not set: %v", nfc.Type)
			}
		})
	}
}

func TestFeatureCollection_MarshalBSON_bbox(t *testing.T) {
	cases := []struct {
		name string
		bbox BBox
	}{
		{
			name: "nil",
			bbox: nil,
		},
		{
			name: "empty",
			bbox: BBox{},
		},
		{
			name: "set",
			bbox: BBox{1, 2, 3, 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fc := NewFeatureCollection()
			fc.BBox = tc.bbox

			data, err := bson.Marshal(fc)
			if err != nil {
				t.Fatalf("unable to marshal feature collection: %e", err)
			}

			nfc := NewFeatureCollection()
			if err = bson.Unmarshal(data, &nfc); err != nil {
				t.Fatalf("unable to unmarshal feature collection: %e", err)
			}

			if !reflect.DeepEqual(nfc.BBox, tc.bbox) {
				t.Errorf("incorrect bbox: %v != %v", nfc.BBox, tc.bbox)
			}
		})
	}
}

func TestFeatureCollection_MarshalBSON_ring(t *testing.T) {
	ring := geo.Ring{{1, 2}, {3, 4}, {5, 6}}
	fc := NewFeatureCollection()
	fc.Append(NewFeature(ring))
	data, err := bson.Marshal(fc)
	if err != nil {
		t.Fatalf("unable to marshal feature collection: %e", err)
	}

	nfc := NewFeatureCollection()
	if err = bson.Unmarshal(data, &nfc); err != nil {
		t.Fatalf("unable to unmarshal feature collection: %e", err)
	}

	if nfc.Features[0].Geometry.GeoJSONType() != "Polygon" {
		t.Errorf("incorrect geometry type: %v != Polygon", nfc.Features[0].Geometry.GeoJSONType())
	}

	if !geo.Equal(nfc.Features[0].Geometry, geo.Polygon{ring}) {
		t.Errorf("incorrect geometry: %v", nfc.Features[0].Geometry)
	}
}

func TestFeatureCollection_MarshalBSON_bound(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{1, 2}, Max: geo.Point{3, 4}}
	fc := NewFeatureCollection()
	fc.Append(NewFeature(bound))
	data, err := bson.Marshal(fc)
	if err != nil {
		t.Fatalf("unable to marshal feature collection: %e", err)
	}

	nfc := NewFeatureCollection()
	if err = bson.Unmarshal(data, &nfc); err != nil {
		t.Fatalf("unable to unmarshal feature collection: %e", err)
	}

	if nfc.Features[0].Geometry.GeoJSONType() != "Polygon" {
		t.Errorf("incorrect geometry type: %v != Polygon", nfc.Features[0].Geometry.GeoJSONType())
	}

	if !geo.Equal(nfc.Features[0].Geometry, bound.ToPolygon()) {
		t.Errorf("incorrect geometry: %v", nfc.Features[0].Geometry)
	}
}

func TestFeatureCollection_MarshalBSON_extraMembers(t *testing.T) {
	fc := NewFeatureCollection()
	fc.Append(NewFeature(geo.Point{1, 2}))
	fc.ExtraMembers = map[string]interface{}{
		"a": 1.0,
		"b": 2.0,
	}

	data, err := bson.Marshal(fc)
	if err != nil {
		t.Fatalf("unable to marshal feature collection: %e", err)
	}

	nfc := NewFeatureCollection()
	if err = bson.Unmarshal(data, &nfc); err != nil {
		t.Fatalf("unable to unmarshal feature collection: %e", err)
	}

	if v := nfc.ExtraMembers["a"]; v != 1.0 {
		t.Errorf("incorrect extra member: %v != %v", v, 1.0)
	}

	if v := nfc.ExtraMembers["b"]; v != 2.0 {
		t.Errorf("incorrect extra member: %v != %v", v, 2.0)
	}
}

func TestUnmarshalFeatureCollection(t *testing.T) {
	rawJSON := `
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      },
	      { "type": "Feature",
	        "geometry": {
	          "type": "LineString",
	          "coordinates": [
	            [102.0, 0.0], [103.0, 1.0], [104.0, 0.0], [105.0, 1.0]
	            ]
	          },
	        "properties": {
	          "prop0": "value0",
	          "prop1": 0.0
	        }
	      },
	      { "type": "Feature",
	         "geometry": {
	           "type": "Polygon",
	           "coordinates": [
	             [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0],
	               [100.0, 1.0], [100.0, 0.0] ]
	             ]
	         },
	         "properties": {
	           "prop0": "value0",
	           "prop1": {"this": "that"}
	         }
	       }
	     ]
	  }`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature collection without issue, err %e", err)
	}

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}

	if len(fc.Features) != 3 {
		t.Errorf("should have 3 features but got %d", len(fc.Features))
	}

	f := fc.Features[0]
	if gt := f.Geometry.GeoJSONType(); gt != "Point" {
		t.Errorf("incorrect feature type: %v != %v", gt, "Point")
	}

	f = fc.Features[1]
	if gt := f.Geometry.GeoJSONType(); gt != "LineString" {
		t.Errorf("incorrect feature type: %v != %v", gt, "LineString")
	}

	f = fc.Features[2]
	if gt := f.Geometry.GeoJSONType(); gt != "Polygon" {
		t.Errorf("incorrect feature type: %v != %v", gt, "Polygon")
	}

	// check unmarshal/marshal loop
	var expected interface{}
	if err = unmarshalJSON([]byte(rawJSON), &expected); err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	data, err := json.MarshalIndent(fc, "", " ")
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	var raw interface{}
	if err = unmarshalJSON(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	if !reflect.DeepEqual(raw, expected) {
		t.Errorf("invalid marshalling: \n%v", string(data))
	}

	// not a feature collection
	data, _ = NewFeature(geo.Point{}).MarshalJSON()
	if _, err = UnmarshalFeatureCollection(data); err == nil {
		t.Error("should return error if not a feature collection")
	}

	if !strings.Contains(err.Error(), "not a feature collection") {
		t.Errorf("incorrect error: %e", err)
	}

	// invalid json
	// truncated
	if _, err = UnmarshalFeatureCollection([]byte(`{"type": "FeatureCollection",`)); err == nil {
		t.Errorf("should return error for invalid json")
	}
}

func TestUnmarshalFeatureCollection_errors(t *testing.T) {
	t.Run("type not a string", func(t *testing.T) {
		rawJSON := `
		  { "type": { "foo":"bar" },
		    "features": [
		      { "type": "Feature",
		        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		        "properties": {"prop0": "value0"}
		      }
		     ]
		  }`

		_, err := UnmarshalFeatureCollection([]byte(rawJSON))
		if _, ok := err.(*json.UnmarshalTypeError); !ok {
			t.Fatalf("wrong error: %T: %e", err, err)
		}
	})

	t.Run("bbox invalid", func(t *testing.T) {
		rawJSON := `
		  { "type": "FeatureCollection",
		    "bbox": { "foo":"bar" },
		    "features": [
		      { "type": "Feature",
		        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		        "properties": {"prop0": "value0"}
		      }
		     ]
		  }`

		_, err := UnmarshalFeatureCollection([]byte(rawJSON))
		if _, ok := err.(*json.UnmarshalTypeError); !ok {
			t.Fatalf("wrong error: %T: %e", err, err)
		}
	})

	t.Run("features invalid", func(t *testing.T) {
		rawJSON := `
		  { "type": "FeatureCollection",
		    "features": { "foo":"bar" }
		  }`

		_, err := UnmarshalFeatureCollection([]byte(rawJSON))
		if _, ok := err.(*json.UnmarshalTypeError); !ok {
			t.Fatalf("wrong error: %T: %e", err, err)
		}
	})
}
