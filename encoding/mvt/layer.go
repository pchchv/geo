package mvt

import "github.com/pchchv/geo/geojson"

const DefaultExtent = 4096

// Layer is intermediate MVT layer to be encoded/decoded or projected.
type Layer struct {
	Name     string
	Version  uint32
	Extent   uint32
	Features []*geojson.Feature
}

// NewLayer is a helper to create a Layer from a feature collection and a name,
// it sets the default extent and version to 1.
func NewLayer(name string, fc *geojson.FeatureCollection) *Layer {
	return &Layer{
		Name:     name,
		Version:  1,
		Extent:   DefaultExtent,
		Features: fc.Features,
	}
}

// Layers is a set of layers.
type Layers []*Layer

// NewLayers creates a set of layers given a set of feature collections.
func NewLayers(layers map[string]*geojson.FeatureCollection) Layers {
	result := make(Layers, 0, len(layers))
	for name, fc := range layers {
		result = append(result, NewLayer(name, fc))
	}

	return result
}

// ToFeatureCollections converts the layers to sets of geojson feature collections.
func (ls Layers) ToFeatureCollections() map[string]*geojson.FeatureCollection {
	result := make(map[string]*geojson.FeatureCollection, len(ls))
	for _, l := range ls {
		result[l.Name] = &geojson.FeatureCollection{
			Features: l.Features,
		}
	}

	return result
}
