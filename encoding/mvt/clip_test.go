package mvt

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
)

func TestLayersClip(t *testing.T) {
	cases := []struct {
		name   string
		bound  geo.Bound
		input  Layers
		output Layers
	}{
		{
			name: "clips polygon and line",
			input: Layers{&Layer{
				Features: []*geojson.Feature{
					geojson.NewFeature(geo.Polygon([]geo.Ring{
						{
							{-10, 10}, {0, 10}, {10, 10}, {10, 5}, {10, -5},
							{10, -10}, {20, -10}, {20, 10}, {40, 10}, {40, 20},
							{20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20},
							{-10, 20},
						},
					})),
					geojson.NewFeature(geo.LineString{{-15, 0}, {66, 0}}),
				},
			}},
			output: Layers{&Layer{
				Features: []*geojson.Feature{
					geojson.NewFeature(geo.Polygon([]geo.Ring{
						{
							{0, 10}, {0, 10}, {10, 10}, {10, 5}, {10, 0},
							{20, 0}, {20, 10}, {30, 10}, {30, 20}, {20, 20},
							{20, 30}, {10, 30}, {10, 20}, {5, 20}, {0, 20},
						},
					})),
					geojson.NewFeature(geo.LineString{{0, 0}, {30, 0}}),
				},
			}},
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{30, 30}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Clip(tc.bound)
			if !reflect.DeepEqual(tc.input, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", tc.input)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestLayerClip_empty(t *testing.T) {
	layer := &Layer{
		Features: []*geojson.Feature{
			geojson.NewFeature(geo.Polygon{{
				{-1, 1}, {0, 1}, {1, 1}, {1, 5}, {1, -5},
			}}),
			geojson.NewFeature(geo.LineString{{55, 0}, {66, 0}}),
		},
	}

	layer.Clip(geo.Bound{Min: geo.Point{50, -10}, Max: geo.Point{70, 10}})
	if v := len(layer.Features); v != 1 {
		t.Errorf("incorrect number of features: %d", v)
	}

	if v := layer.Features[0].Geometry.GeoJSONType(); v != "LineString" {
		t.Errorf("kept the wrong geometry: %v", layer.Features[0].Geometry)
	}
}
