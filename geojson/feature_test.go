package geojson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func TestNewFeature(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if f.Type != "Feature" {
		t.Errorf("incorrect feature: %v != Feature", f.Type)
	}
}

func TestFeatureMarshalJSON(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshalJSON_BBox(t *testing.T) {
	f := NewFeature(geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{2, 2}})
	f.BBox = nil
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	} else if bytes.Contains(blob, []byte(`"bbox"`)) {
		t.Errorf("should not set the bbox value")
	}

	f.BBox = []float64{1, 2, 3, 4}
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	} else if !bytes.Contains(blob, []byte(`"bbox":[1,2,3,4]`)) {
		t.Errorf("should set type to polygon coords: %v", string(blob))
	}
}

func TestFeatureMarshalJSON_Bound(t *testing.T) {
	f := NewFeature(geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{2, 2}})
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	} else if !bytes.Contains(blob, []byte(`"type":"Polygon"`)) {
		t.Errorf("should set type to polygon")
	} else if !bytes.Contains(blob, []byte(`"coordinates":[[[1,1],[2,1],[2,2],[1,2],[1,1]]]`)) {
		t.Errorf("should set type to polygon coords: %v", string(blob))
	}
}

func TestFeature_marshalValue(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if blob, err := json.Marshal(*f); err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}

	if blob, err := bson.Marshal(*f); err != nil {
		t.Fatalf("should marshal to bson just fine but got %v", err)
	} else if !bytes.Contains(blob, append([]byte{byte(bsontype.Null)}, []byte("properties")...)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshal(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if blob, err := json.Marshal(f); err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	} else if !bytes.Contains(blob, []byte(`"type":"Feature"`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestUnmarshalFeature_GeometryCollection(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type":"GeometryCollection","geometries":[{"type": "Point", "coordinates": [102.0, 0.5]}]}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	wantType := geo.Collection{}.GeoJSONType()
	if f.Geometry.GeoJSONType() != wantType {
		t.Fatalf("invalid GeoJSONType: %v", f.Geometry.GeoJSONType())
	}
}

