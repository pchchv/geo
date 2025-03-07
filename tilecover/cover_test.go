package tilecover

import (
	"encoding/json"
	"math"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/planar"
)

func TestCountries(t *testing.T) {
	files, err := os.ReadDir("./testdata/world")
	if err != nil {
		t.Errorf("could not read directory: %e", err)
	}

	var countries []string
	for _, info := range files {
		if !strings.Contains(info.Name(), "_out") {
			countries = append(countries, strings.Split(info.Name(), ".")[0])
		}
	}

	for _, country := range countries {
		t.Run(country, func(t *testing.T) {
			f := loadFeature(t, "./testdata/world/"+country+".geo.json")
			tiles, _ := Geometry(f.Geometry, 6)
			tiles = MergeUp(tiles, 1)
			expected := loadFeatureCollection(t, "./testdata/world/"+country+"_out.geojson")
			compareFeatureCollections(t, country, tiles.ToFeatureCollection(), expected)
		})
	}
}

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

func compareFeatureCollections(t testing.TB, name string, result, expected *geojson.FeatureCollection) {
	sortFC(result)
	sortFC(expected)

	t.Helper()

	if len(result.Features) != len(expected.Features) {
		t.Errorf("feature count mismatch: %v != %v", len(result.Features), len(expected.Features))
		output(t, name, result)
		return
	}

	failure := false
	for i := range result.Features {
		r := result.Features[i].Geometry.(geo.Polygon)
		e := expected.Features[i].Geometry.(geo.Polygon)
		rc, ra := planar.CentroidArea(r)
		ec, ea := planar.CentroidArea(e)

		if delta := math.Abs(ra - ea); delta > 0.01 {
			failure = true
			t.Errorf("f %d: area not equal: %v", i, delta)
		}

		if dist := geo.Distance(rc, ec); dist > 1 {
			failure = true
			t.Errorf("f %d: centroid far apart: %v", i, dist)
		}
	}

	if failure {
		output(t, name, result)
	}
}
