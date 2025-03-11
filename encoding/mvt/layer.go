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
