package tilecover

import (
	"os"
	"testing"

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
