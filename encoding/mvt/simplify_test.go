package mvt

import (
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/simplifier"
)

func TestLayerSimplify(t *testing.T) {
	// should remove feature that are empty.
	ls := Layers{&Layer{
		Features: []*geojson.Feature{
			geojson.NewFeature(geo.LineString(nil)),
			geojson.NewFeature(geo.LineString{{0, 0}, {1, 1}}),
		},
	}}

	simplifier := simplifier.DouglasPeucker(10)
	ls.Simplify(simplifier)
	if len(ls[0].Features) != 1 {
		t.Errorf("should remove empty feature")
	}

	if v := ls[0].Features[0].Geometry.GeoJSONType(); v != "LineString" {
		t.Errorf("incorrect type: %v", v)
	}
}
