package geojson

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
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
		t.Fatalf("error marshalling to json: %e", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshalJSON_BBox(t *testing.T) {
	f := NewFeature(geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{2, 2}})
	f.BBox = nil
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %e", err)
	} else if bytes.Contains(blob, []byte(`"bbox"`)) {
		t.Errorf("should not set the bbox value")
	}

	f.BBox = []float64{1, 2, 3, 4}
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %e", err)
	} else if !bytes.Contains(blob, []byte(`"bbox":[1,2,3,4]`)) {
		t.Errorf("should set type to polygon coords: %v", string(blob))
	}
}

func TestFeatureMarshalJSON_Bound(t *testing.T) {
	f := NewFeature(geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{2, 2}})
	if blob, err := f.MarshalJSON(); err != nil {
		t.Fatalf("error marshalling to json: %e", err)
	} else if !bytes.Contains(blob, []byte(`"type":"Polygon"`)) {
		t.Errorf("should set type to polygon")
	} else if !bytes.Contains(blob, []byte(`"coordinates":[[[1,1],[2,1],[2,2],[1,2],[1,1]]]`)) {
		t.Errorf("should set type to polygon coords: %v", string(blob))
	}
}

func TestFeature_marshalValue(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if blob, err := json.Marshal(*f); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}

	if blob, err := bson.Marshal(*f); err != nil {
		t.Fatalf("should marshal to bson just fine but got %e", err)
	} else if !bytes.Contains(blob, append([]byte{byte(bsontype.Null)}, []byte("properties")...)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshal(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if blob, err := json.Marshal(f); err != nil {
		t.Fatalf("should marshal to json just fine but got %e", err)
	} else if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	} else if !bytes.Contains(blob, []byte(`"type":"Feature"`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestUnmarshalFeature(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	    "properties": {"prop0": "value0"}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	if f.Type != "Feature" {
		t.Errorf("should have type of Feature got: %v", f.Type)
	}

	if len(f.Properties) != 1 {
		t.Errorf("should have 1 property but got: %v", f.Properties)
	}

	// not a feature
	data, _ := NewFeatureCollection().MarshalJSON()
	if _, err = UnmarshalFeature(data); err == nil {
		t.Error("should return error if not a feature")
	} else if !strings.Contains(err.Error(), "not a feature") {
		t.Errorf("incorrect error: %e", err)
	}

	// invalid json
	// truncated
	if _, err = UnmarshalFeature([]byte(`{"type": "Feature",`)); err == nil {
		t.Errorf("should return error for invalid json")
	}

	f = &Feature{}
	// truncated
	if err = f.UnmarshalJSON([]byte(`{"type": "Feature",`)); err == nil {
		t.Errorf("should return error for invalid json")
	}
}

func TestUnmarshalFeature_GeometryCollection(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type":"GeometryCollection","geometries":[{"type": "Point", "coordinates": [102.0, 0.5]}]}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	wantType := geo.Collection{}.GeoJSONType()
	if f.Geometry.GeoJSONType() != wantType {
		t.Fatalf("invalid GeoJSONType: %v", f.Geometry.GeoJSONType())
	}
}

func TestUnmarshalFeature_missingGeometry(t *testing.T) {
	t.Run("empty geometry", func(t *testing.T) {
		rawJSON := `{ "type": "Feature", "geometry": {} }`
		if _, err := UnmarshalFeature([]byte(rawJSON)); err != ErrInvalidGeometry {
			t.Fatalf("incorrect unmarshal error: %e", err)
		}
	})

	t.Run("missing geometry", func(t *testing.T) {
		rawJSON := `{ "type": "Feature" }`
		if f, err := UnmarshalFeature([]byte(rawJSON)); err != nil {
			t.Fatalf("should not error: %e", err)
		} else if f == nil {
			t.Fatalf("feature should not be nil")
		}
	})
}

func TestFeatureMarshalJSON_null(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		type S struct {
			GeoJSON *Feature `json:"geojson"`
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
			GeoJSON Feature `json:"geojson"`
		}

		var s S
		if err := json.Unmarshal([]byte(`{"geojson": null}`), &s); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		}

		if !reflect.DeepEqual(s.GeoJSON, Feature{}) {
			t.Errorf("should be empty, got: %v", s)
		}
	})
}

func TestUnmarshalBSON_missingGeometry(t *testing.T) {
	t.Run("missing geometry", func(t *testing.T) {
		f := NewFeature(nil)
		f.Geometry = nil

		data, err := bson.Marshal(f)
		if err != nil {
			t.Fatalf("marshal error: %e", err)
		}

		nf := &Feature{}
		if err = bson.Unmarshal(data, &nf); err != nil {
			t.Fatalf("unmarshal error: %e", err)
		}

		if f.Geometry != nil {
			t.Fatalf("geometry should be nil")
		}

		if f == nil {
			t.Fatalf("feature should not be nil")
		}
	})
}

func TestUnmarshalFeature_BBox(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		"bbox": [1,2,3,4],
	    "properties": {"prop0": "value0"}
	  }`

	if f, err := UnmarshalFeature([]byte(rawJSON)); err != nil {
		t.Fatalf("unmarshal error: %e", err)
	} else if !f.BBox.Valid() {
		t.Errorf("bbox should be valid: %v", f.BBox)
	}
}

func TestMarshalFeatureID(t *testing.T) {
	f := &Feature{
		ID: "asdf",
	}

	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %e", err)
	} else if !bytes.Equal(data, []byte(`{"id":"asdf","type":"Feature","geometry":null,"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}

	f.ID = 123
	if data, err = f.MarshalJSON(); err != nil {
		t.Fatalf("should marshal, %e", err)
	} else if !bytes.Equal(data, []byte(`{"id":123,"type":"Feature","geometry":null,"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}

func TestUnmarshalFeatureID(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "id": 123,
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature without issue, err %e", err)
	}

	if v, ok := f.ID.(float64); !ok || v != 123 {
		t.Errorf("should parse id as number, got %T %f", f.ID, v)
	}

	rawJSON = `
	  { "type": "Feature",
	    "id": "abcd",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]}
	  }`

	f, err = UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature without issue, err %e", err)
	}

	if v, ok := f.ID.(string); !ok || v != "abcd" {
		t.Errorf("should parse id as string, got %T %s", f.ID, v)
	}
}

func TestMarshalRing(t *testing.T) {
	ring := geo.Ring{{0, 0}, {1, 1}, {2, 1}, {0, 0}}
	f := NewFeature(ring)
	if data, err := f.MarshalJSON(); err != nil {
		t.Fatalf("should marshal, %e", err)
	} else if !bytes.Equal(data, []byte(`{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[0,0],[1,1],[2,1],[0,0]]]},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}

// // uncomment to test/benchmark custom json marshalling
// func init() {
// 	var c = jsoniter.Config{
// 		EscapeHTML:              true,
// 		SortMapKeys:             false,
// 		ValidateJsonRawMessage:  false,
// 		MarshalFloatWith6Digits: true,
// 	}.Froze()

// 	CustomJSONMarshaler = c
// 	CustomJSONUnmarshaler = c
// }

func BenchmarkFeatureMarshalJSON(b *testing.B) {
	data, err := os.ReadFile("../encoding/mvt/testdata/16-17896-24449.json")
	if err != nil {
		b.Fatalf("could not open file: %e", err)
	}

	tile := map[string]*FeatureCollection{}
	if err = json.Unmarshal(data, &tile); err != nil {
		b.Fatalf("could not unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := marshalJSON(tile)
		if err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkFeatureUnmarshalJSON(b *testing.B) {
	data, err := os.ReadFile("../encoding/mvt/testdata/16-17896-24449.json")
	if err != nil {
		b.Fatalf("could not open file: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tile := map[string]*FeatureCollection{}
		if err = unmarshalJSON(data, &tile); err != nil {
			b.Fatalf("could not unmarshal: %e", err)
		}
	}
}

func BenchmarkFeatureMarshalBSON(b *testing.B) {
	data, err := os.ReadFile("../encoding/mvt/testdata/16-17896-24449.json")
	if err != nil {
		b.Fatalf("could not open file: %e", err)
	}

	tile := map[string]*FeatureCollection{}
	if err = json.Unmarshal(data, &tile); err != nil {
		b.Fatalf("could not unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bson.Marshal(tile); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkFeatureUnmarshalBSON(b *testing.B) {
	data, err := os.ReadFile("../encoding/mvt/testdata/16-17896-24449.json")
	if err != nil {
		b.Fatalf("could not open file: %e", err)
	}

	tile := map[string]*FeatureCollection{}
	if err = json.Unmarshal(data, &tile); err != nil {
		b.Fatalf("could not unmarshal: %e", err)
	}

	bdata, err := bson.Marshal(tile)
	if err != nil {
		b.Fatalf("could not marshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tile := map[string]*FeatureCollection{}
		if err = bson.Unmarshal(bdata, &tile); err != nil {
			b.Fatalf("could not unmarshal: %e", err)
		}
	}
}
