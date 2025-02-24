package geojson

import (
	"bytes"
	"testing"

	"github.com/pchchv/geo"
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
