package geojson

import "github.com/pchchv/geo"

type featureDoc struct {
	ID         interface{} `json:"id,omitempty" bson:"id"`
	Type       string      `json:"type" bson:"type"`
	BBox       BBox        `json:"bbox,omitempty" bson:"bbox,omitempty"`
	Geometry   *Geometry   `json:"geometry" bson:"geometry"`
	Properties Properties  `json:"properties" bson:"properties"`
}

func newFeatureDoc(f *Feature) *featureDoc {
	doc := &featureDoc{
		ID:         f.ID,
		Type:       "Feature",
		Properties: f.Properties,
		BBox:       f.BBox,
		Geometry:   NewGeometry(f.Geometry),
	}
	if len(doc.Properties) == 0 {
		doc.Properties = nil
	}

	return doc
}

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
