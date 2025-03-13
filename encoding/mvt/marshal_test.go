package mvt

import (
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
