package geojson

import "go.mongodb.org/mongo-driver/bson"

const featureCollection = "FeatureCollection"

// A FeatureCollection correlates to a GeoJSON feature collection.
type FeatureCollection struct {
	// ExtraMembers can be used to encoded/decode extra key/members in the
	// base of the feature collection. Note that keys of "type",
	// "bbox" and "features" will not work as those are reserved by the GeoJSON spec.
	ExtraMembers Properties `json:"-"`
	Type         string     `json:"type"`
	BBox         BBox       `json:"bbox,omitempty"`
	Features     []*Feature `json:"features"`
}

// NewFeatureCollection creates and initializes a new feature collection.
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     featureCollection,
		Features: []*Feature{},
	}
}

// MarshalJSON converts the feature collection object into the proper JSON.
// It will handle the encoding of all the child features and geometries.
// Alternately one can call json.Marshal(fc) directly for the same result.
// Items in the ExtraMembers map will be included in the base of the feature collection object.
func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	m := newFeatureCollectionDoc(fc)
	return marshalJSON(m)
}

// MarshalBSON converts the feature collection object
// into a BSON document represented by bytes.
// It will handle the encoding of all the child features and geometries.
// Items in the ExtraMembers map will be included in the
// base of the feature collection object.
func (fc FeatureCollection) MarshalBSON() ([]byte, error) {
	m := newFeatureCollectionDoc(fc)
	return bson.Marshal(m)

func newFeatureCollectionDoc(fc FeatureCollection) (temp map[string]interface{}) {
	if fc.ExtraMembers != nil {
		temp = fc.ExtraMembers.Clone()
	} else {
		temp = make(map[string]interface{}, 3)
	}

	temp["type"] = featureCollection
	delete(temp, "bbox")
	if fc.BBox != nil {
		temp["bbox"] = fc.BBox
	}

	if fc.Features == nil {
		temp["features"] = []*Feature{}
	} else {
		temp["features"] = fc.Features
	}

	return temp
}
