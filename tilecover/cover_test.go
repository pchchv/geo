package tilecover

import (
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
)

func loadFeature(t testing.TB, path string) *geojson.Feature {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to read file: %e", err)
	}

	if f, err := geojson.UnmarshalFeature(data); err == nil {
		return f
	}

	if fc, err := geojson.UnmarshalFeatureCollection(data); err == nil {
		if len(fc.Features) != 1 {
			t.Fatalf("must have 1 feature: %v", len(fc.Features))
		}
		return fc.Features[0]
	}

	g, err := geojson.UnmarshalGeometry(data)
	if err != nil {
		t.Fatalf("unable to unmarshal feature: %v", err)
	}

	return geojson.NewFeature(g.Geometry())
}

func loadFeatureCollection(t testing.TB, path string) *geojson.FeatureCollection {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to read file: %v", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		t.Fatalf("unable to unmarshal feature: %v", err)
	}

	var count int
	for i := range fc.Features {
		if fc.Features[i].Properties["name"] != "original" {
			fc.Features[count] = fc.Features[i]
			count++
		}
	}

	fc.Features = fc.Features[:count]
	return fc
}

func sortFC(fc *geojson.FeatureCollection) {
	sort.Slice(fc.Features, func(i, j int) bool {
		a := fc.Features[i].Geometry.(geo.Polygon)[0]
		b := fc.Features[j].Geometry.(geo.Polygon)[0]
		if a[0][0] != b[0][0] {
			return a[0][0] < b[0][0]
		}

		return a[0][1] < b[0][1]
	})
}

// output gets called if there is a test failure for debugging.
func output(t testing.TB, name string, r *geojson.FeatureCollection) {
	f := loadFeature(t, "./testdata/"+name+".geojson")
	if f.Properties == nil {
		f.Properties = make(geojson.Properties)
	}

	f.Properties["fill"] = "#FF0000"
	f.Properties["fill-opacity"] = "0.5"
	f.Properties["stroke"] = "#FF0000"
	f.Properties["name"] = "original"
	r.Append(f)
	data, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		t.Fatalf("error marshalling json: %v", err)
	}

	if err = os.WriteFile("failure_"+name+".geojson", data, 0644); err != nil {
		t.Fatalf("write file failure: %v", err)
	}
}
