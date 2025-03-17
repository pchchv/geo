package mvt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/maptile"
)

func TestMarshalUnmarshalGzipped_Full(t *testing.T) {
	tile := maptile.New(8956, 12223, 15)
	ls := geo.LineString{
		{-81.60346275, 41.50998572},
		{-81.6033669, 41.50991259},
		{-81.60355599, 41.50976036},
		{-81.6040648, 41.50936811},
		{-81.60404411, 41.50935405},
	}
	expected := ls.Clone()

	f := geojson.NewFeature(ls)
	f.Properties = geojson.Properties{
		"source":       "openstreetmap.org",
		"kind":         "path",
		"name":         "Uptown Alley",
		"landuse_kind": "retail",
		"sort_rank":    float64(354),
		"kind_detail":  "pedestrian",
		"min_zoom":     float64(13),
		"id":           float64(246698394),
	}

	fc := geojson.NewFeatureCollection()
	fc.Append(f)

	layers := Layers{NewLayer("roads", fc)}

	// project to the tile coords
	layers.ProjectToTile(tile)

	// marshal
	encoded, err := MarshalGzipped(layers)
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	}

	// unmarshal
	decoded, err := UnmarshalGzipped(encoded)
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	// project back
	decoded.ProjectToWGS84(tile)

	// compare the results
	result := decoded[0].Features[0]
	compareProperties(t, result.Properties, f.Properties)

	// compare geometry
	xe, ye := tileEpsilon(tile)
	compareGEOGeometry(t, result.Geometry, expected, xe, ye)
}

func TestMarshalUnmarshalForGeometryCollection(t *testing.T) {
	tile := maptile.New(8956, 12223, 15)
	outerRing := geo.Ring{
		{-81.6033000, 41.5099000},
		{-81.6033000, 41.5094000},
		{-81.6039000, 41.5094000},
		{-81.6039000, 41.5099000},
		{-81.6033000, 41.5099000},
	}
	hole := geo.Ring{
		{-81.60389989614487, 41.50941085679876},
		{-81.60329908132553, 41.50941085679876},
		{-81.60329908132553, 41.50990496161759},
		{-81.60389989614487, 41.50990496161759},
		{-81.60389989614487, 41.50941085679876},
	}
	outerPolygon := geo.Polygon{outerRing, hole}
	polygonInHole := geo.Polygon{geo.Ring{
		{-81.60375505685806, 41.5095494475553},
		{-81.6034385561943, 41.5095494475553},
		{-81.6034385561943, 41.50978043149024},
		{-81.60375505685806, 41.50978043149024},
		{-81.60375505685806, 41.5095494475553},
	}}
	geometryCollection := geo.Collection{outerPolygon, polygonInHole}
	expected := []geo.Geometry{outerPolygon.Clone(), outerPolygon.Clone()}

	f := geojson.NewFeature(geometryCollection)
	f.Properties = geojson.Properties{
		"id": float64(246698394),
	}

	fc := geojson.NewFeatureCollection()
	fc.Append(f)
	layers := Layers{NewLayer("roads", fc)}

	// project to the tile coords
	layers.ProjectToTile(tile)

	// marshal
	encoded, err := MarshalGzipped(layers)
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	}

	// unmarshal
	decoded, err := UnmarshalGzipped(encoded)
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	// project back
	decoded.ProjectToWGS84(tile)

	// compare the results
	results := decoded[0].Features
	compareProperties(t, results[0].Properties, f.Properties)

	// compare geometry
	xe, ye := tileEpsilon(tile)
	if len(results) == len(expected) {
		t.Errorf("result geometry count must be splited polygon: %v (but result is %v)", len(results), len(expected))
	}

	for i, result := range results {
		compareGEOGeometry(t, result.Geometry, expected[i], xe, ye)
	}
}

func TestMarshal_nilGeometry(t *testing.T) {
	// should handle nil geometry,
	// possibly caused by clipping or other issues.
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(nil))
	layers := Layers{NewLayer("roads", fc)}

	// project to the tile coords
	tile := maptile.New(0, 0, 15)
	layers.ProjectToTile(tile)

	// marshal
	encoded, err := MarshalGzipped(layers)
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	}

	// unmarshal
	decoded, err := UnmarshalGzipped(encoded)
	if err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	// project back
	decoded.ProjectToWGS84(tile)
	if v := len(decoded[0].Features); v != 0 {
		t.Errorf("should have have any features, has %v", v)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		tile maptile.Tile
	}{
		{
			name: "15-8956-12223",
			tile: maptile.New(8956, 12223, 15),
		},
		{
			name: "16-17896-24449",
			tile: maptile.New(17896, 24449, 16),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := loadGeoJSON(t, tc.tile)
			layers := NewLayers(loadGeoJSON(t, tc.tile))
			layers.ProjectToTile(tc.tile)
			data, err := Marshal(layers)
			if err != nil {
				t.Errorf("marshal error: %e", err)
			}

			layers, err = Unmarshal(data)
			if err != nil {
				t.Errorf("unmarshal error: %e", err)
			}

			layers.ProjectToWGS84(tc.tile)
			result := layers.ToFeatureCollections()
			xEpsilon, yEpsilon := tileEpsilon(tc.tile)
			for key := range expected {
				for i := range expected[key].Features {
					r := result[key].Features[i]
					e := expected[key].Features[i]

					// t.Logf("checking layer %s: feature %d", key, i)
					compareProperties(t, r.Properties, e.Properties)
					compareGEOGeometry(t, r.Geometry, e.Geometry, xEpsilon, yEpsilon)
				}
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		tile maptile.Tile
	}{
		{
			name: "15-8956-12223",
			tile: maptile.New(8956, 12223, 15),
		},
		{
			name: "16-17896-24449",
			tile: maptile.New(17896, 24449, 16),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := loadGeoJSON(t, tc.tile)
			layers, err := UnmarshalGzipped(loadMVT(t, tc.tile))
			if err != nil {
				t.Fatalf("error unmarshalling: %e", err)
			}

			layers.ProjectToWGS84(tc.tile)
			result := layers.ToFeatureCollections()
			xEpsilon, yEpsilon := tileEpsilon(tc.tile)
			for key := range expected {
				for i := range expected[key].Features {
					r := result[key].Features[i]
					e := expected[key].Features[i]
					t.Logf("checking layer %s: feature %d", key, i)
					compareProperties(t, r.Properties, e.Properties)
					compareGEOGeometry(t, r.Geometry, e.Geometry, xEpsilon, yEpsilon)
				}
			}
		})
	}
}

func TestUnmarshalGzippedTileWithRegularUnmarshalFunction(t *testing.T) {
	t.Run("15-8956-12223", func(t *testing.T) {
		tile := maptile.New(8956, 12223, 15)
		if _, err := Unmarshal(loadMVT(t, tile)); err.Error() != "failed to unmarshal, data possibly gzipped" {
			t.Fatal()
		}
	})
}

func TestMarshal_ID(t *testing.T) {
	cases := []struct {
		name string
		id   interface{}
		val  float64
	}{
		{
			name: "int",
			id:   int(86427531),
			val:  86427531,
		},
		{
			name: "int8",
			id:   int8(123),
			val:  123,
		},
		{
			name: "int16",
			id:   int16(6884),
			val:  6884,
		},
		{
			name: "int32",
			id:   int32(123),
			val:  123,
		},
		{
			name: "int64",
			id:   int64(12345678),
			val:  12345678,
		},
		{
			name: "uint",
			id:   uint(86427531),
			val:  86427531,
		},
		{
			name: "uint8",
			id:   uint8(123),
			val:  123,
		},
		{
			name: "uint16",
			id:   uint16(6884),
			val:  6884,
		},
		{
			name: "uint32",
			id:   uint32(123),
			val:  123,
		},
		{
			name: "uint64",
			id:   uint64(12345678),
			val:  12345678,
		},
		{
			name: "float32",
			id:   float32(123.45),
			val:  123,
		},
		{
			name: "float64",
			id:   float64(123.45),
			val:  123,
		},
		{
			name: "string",
			id:   "123456",
			val:  123456,
		},

		// negatives
		{
			name: "negative string",
			id:   "-123456",
			val:  0, // nil
		},
		{
			name: "negative int",
			id:   int(-123456),
			val:  0, // nil
		},
		{
			name: "negative int64",
			id:   int64(-123456),
			val:  0, // nil
		},
		{
			name: "negative float64",
			id:   float64(-123456),
			val:  0, // nil
		},
	}

	f := geojson.NewFeature(geo.Point{1, 2})
	f.Properties["type"] = "point"
	fc := geojson.NewFeatureCollection().Append(f)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f.ID = tc.id
			data, err := Marshal(NewLayers(map[string]*geojson.FeatureCollection{"roads": fc}))
			if err != nil {
				t.Errorf("marshal error: %e", err)
			}

			ls, err := Unmarshal(data)
			if err != nil {
				t.Errorf("unmarshal error: %e", err)
			}

			id := ls.ToFeatureCollections()["roads"].Features[0].ID
			if tc.val > 0 {
				if id.(float64) != tc.val {
					t.Errorf("incorrect id: %v != %v", id, tc.val)
				}
			} else {
				if id != nil {
					t.Errorf("id should be nil: %v", id)
				}
			}
		})
	}

	t.Run("unmarshaled int from json", func(t *testing.T) {
		f.ID = 123
		data, err := fc.MarshalJSON()
		if err != nil {
			t.Fatalf("json marshal error: %e", err)
		}

		fc, err = geojson.UnmarshalFeatureCollection(data)
		if err != nil {
			t.Fatalf("unmarshal json error: %e", err)
		}

		if _, ok := fc.Features[0].ID.(float64); !ok {
			t.Errorf("json should unmarshal numbers to float64: %T", fc.Features[0].ID)
		}

		data, err = Marshal(NewLayers(map[string]*geojson.FeatureCollection{"roads": fc}))
		if err != nil {
			t.Errorf("marshal error: %e", err)
		}

		ls, err := Unmarshal(data)
		if err != nil {
			t.Errorf("unmarshal error: %e", err)
		}

		id := ls.ToFeatureCollections()["roads"].Features[0].ID
		if _, ok := id.(float64); !ok {
			t.Errorf("should unmarshal id to float64: %T", id)
		}
	})
}

func TestStableMarshalling(t *testing.T) {
	layers := NewLayers(loadGeoJSON(t, maptile.New(17896, 24449, 16)))
	firstData, _ := Marshal(layers)
	for i := 0; i < 100; i++ {
		data, _ := Marshal(layers)
		if !bytes.Equal(data, firstData) {
			t.Errorf("a marshal had different bytes")
		}
	}
}

func compareProperties(t testing.TB, result, expected geojson.Properties) {
	t.Helper()

	// properties
	fr := map[string]interface{}(result)
	fe := map[string]interface{}(expected)
	for k, v := range fe {
		if _, ok := v.([]interface{}); ok {
			// arrays are not included
			delete(fr, k)
			delete(fe, k)
		}

		if k == "scale_rank" || k == "layer" {
			if v == 1.0 {
				delete(fr, k)
				delete(fe, k)
			}
		}
	}

	if !reflect.DeepEqual(fr, fe) {
		t.Errorf("properties not equal")
		if len(fr) != len(fe) {
			t.Errorf("properties length not equal: %v != %v", len(fr), len(fe))
		}

		for k := range fr {
			t.Logf("%s: %T %v -- %T %v", k, fr[k], fr[k], fe[k], fe[k])
		}
	}
}

func compareGEOGeometry(t testing.TB, result, expected geo.Geometry, xEpsilon, yEpsilon float64) {
	t.Helper()
	if result.GeoJSONType() != expected.GeoJSONType() {
		t.Errorf("different types: %v != %v", result.GeoJSONType(), expected.GeoJSONType())
		return
	}

	switch r := result.(type) {
	case geo.Point:
		comparePoints(t,
			[]geo.Point{r},
			[]geo.Point{expected.(geo.Point)},
			xEpsilon, yEpsilon,
		)
	case geo.MultiPoint:
		comparePoints(t,
			[]geo.Point(r),
			[]geo.Point(expected.(geo.MultiPoint)),
			xEpsilon, yEpsilon,
		)
	case geo.LineString:
		comparePoints(t,
			[]geo.Point(r),
			[]geo.Point(expected.(geo.LineString)),
			xEpsilon, yEpsilon,
		)
	case geo.MultiLineString:
		e := expected.(geo.MultiLineString)
		for i := range r {
			compareGEOGeometry(t, r[i], e[i], xEpsilon, yEpsilon)
		}
	case geo.Polygon:
		e := expected.(geo.Polygon)
		for i := range r {
			compareGEOGeometry(t, geo.LineString(r[i]), geo.LineString(e[i]), xEpsilon, yEpsilon)
		}
	case geo.MultiPolygon:
		e := expected.(geo.MultiPolygon)
		for i := range r {
			compareGEOGeometry(t, r[i], e[i], xEpsilon, yEpsilon)
		}
	default:
		t.Errorf("unsupported type: %T", result)
	}
}

func comparePoints(t testing.TB, e, r []geo.Point, xEpsilon, yEpsilon float64) {
	if len(r) != len(e) {
		t.Errorf("geometry length not equal: %v != %v", len(r), len(e))
	}

	for i := range e {
		xe := math.Abs(r[i][0] - e[i][0])
		if xe > xEpsilon {
			t.Errorf("%d x: %f != %f    %f", i, r[i][0], e[i][0], xe)
		}

		ye := math.Abs(r[i][1] - e[i][1])
		if ye > yEpsilon {
			t.Errorf("%d y: %f != %f    %f", i, r[i][1], e[i][1], ye)
		}
	}
}

func loadMVT(t testing.TB, tile maptile.Tile) []byte {
	data, err := os.ReadFile(fmt.Sprintf("testdata/%d-%d-%d.mvt", tile.Z, tile.X, tile.Y))
	if err != nil {
		t.Fatalf("failed to load mvt file: %e", err)
	}

	return data
}

func loadGeoJSON(t testing.TB, tile maptile.Tile) map[string]*geojson.FeatureCollection {
	data, err := os.ReadFile(fmt.Sprintf("testdata/%d-%d-%d.json", tile.Z, tile.X, tile.Y))
	if err != nil {
		t.Fatalf("failed to load mvt file: %e", err)
	}

	r := make(map[string]*geojson.FeatureCollection)
	if err = json.Unmarshal(data, &r); err != nil {
		t.Fatalf("unmarshal error: %e", err)
	}

	return r
}

func tileEpsilon(tile maptile.Tile) (float64, float64) {
	b := tile.Bound()
	xEpsilon := (b.Max[0] - b.Min[0]) / 4096 * 2 // allowed error
	yEpsilon := (b.Max[1] - b.Min[1]) / 4096 * 2

	return xEpsilon, yEpsilon
}

func BenchmarkMarshal(b *testing.B) {
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(layers)
		if err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	data, err := Marshal(layers)
	if err != nil {
		b.Fatalf("marshal error: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(data)
		if err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkProjectToTile(b *testing.B) {
	tile := maptile.New(17896, 24449, 16)
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		layers.ProjectToTile(tile)
	}
}

func BenchmarkProjectToGeo(b *testing.B) {
	tile := maptile.New(17896, 24449, 16)
	layers := NewLayers(loadGeoJSON(b, maptile.New(17896, 24449, 16)))
	layers.ProjectToTile(tile)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		layers.ProjectToWGS84(tile)
	}
}
