package geojson

import (
	"bytes"
	"fmt"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
)

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

// MarshalJSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
func (f Feature) MarshalJSON() ([]byte, error) {
	return marshalJSON(newFeatureDoc(&f))
}

// MarshalBSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
func (f Feature) MarshalBSON() ([]byte, error) {
	return bson.Marshal(newFeatureDoc(&f))
}

// UnmarshalJSON handles the correct unmarshalling of the data
// into the geo.Geometry types.
func (f *Feature) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`null`)) {
		*f = Feature{}
		return nil
	}

	doc := &featureDoc{}
	if err := unmarshalJSON(data, &doc); err != nil {
		return err
	}

	return featureUnmarshalFinish(doc, f)
}

// UnmarshalBSON will unmarshal a BSON document created with bson.Marshal.
func (f *Feature) UnmarshalBSON(data []byte) error {
	doc := &featureDoc{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}

	return featureUnmarshalFinish(doc, f)
}

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

func featureUnmarshalFinish(doc *featureDoc, f *Feature) error {
	if doc.Type != "Feature" {
		return fmt.Errorf("geojson: not a feature: type=%s", doc.Type)
	}

	var g geo.Geometry
	if doc.Geometry != nil {
		if doc.Geometry.Coordinates == nil && doc.Geometry.Geometries == nil {
			return ErrInvalidGeometry
		}
		g = doc.Geometry.Geometry()
	}

	*f = Feature{
		ID:         doc.ID,
		Type:       doc.Type,
		Properties: doc.Properties,
		BBox:       doc.BBox,
		Geometry:   g,
	}

	return nil
}
