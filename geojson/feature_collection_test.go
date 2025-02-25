package geojson

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
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
