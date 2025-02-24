package geojson

import "github.com/pchchv/geo"

// Feature corresponds to GeoJSON feature object.
type Feature struct {
	ID         interface{}  `json:"id,omitempty"`
	Type       string       `json:"type"`
	BBox       BBox         `json:"bbox,omitempty"`
	Geometry   geo.Geometry `json:"geometry"`
	Properties Properties   `json:"properties"`
}

// NewFeature creates and initializes a GeoJSON feature given the required attributes.
func NewFeature(geometry geo.Geometry) *Feature {
	return &Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(map[string]interface{}),
	}
}
