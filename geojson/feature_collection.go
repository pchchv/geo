package geojson

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

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
}

// UnmarshalJSON decodes the data into a GeoJSON feature collection.
// Extra/foreign members will be put into the `ExtraMembers` attribute.
func (fc *FeatureCollection) UnmarshalJSON(data []byte) (err error) {
	if bytes.Equal(data, []byte(`null`)) {
		*fc = FeatureCollection{}
		return nil
	}

	temp := make(map[string]nocopyRawMessage, 4)
	if err = unmarshalJSON(data, &temp); err != nil {
		return
	}

	*fc = FeatureCollection{}
	for key, value := range temp {
		switch key {
		case "type":
			if err = unmarshalJSON(value, &fc.Type); err != nil {
				return
			}
		case "bbox":
			if err = unmarshalJSON(value, &fc.BBox); err != nil {
				return
			}
		case "features":
			if err = unmarshalJSON(value, &fc.Features); err != nil {
				return
			}
		default:
			if fc.ExtraMembers == nil {
				fc.ExtraMembers = Properties{}
			}

			var val interface{}
			if err = unmarshalJSON(value, &val); err != nil {
				return
			} else {
				fc.ExtraMembers[key] = val
			}
		}
	}

	if fc.Type != featureCollection {
		return fmt.Errorf("geojson: not a feature collection: type=%s", fc.Type)
	}

	return nil
}

// UnmarshalBSON will unmarshal a BSON document created with bson.Marshal.
// Extra/foreign members will be put into the `ExtraMembers` attribute.
func (fc *FeatureCollection) UnmarshalBSON(data []byte) (err error) {
	temp := make(map[string]bson.RawValue, 4)

	if err = bson.Unmarshal(data, &temp); err != nil {
		return
	}

	*fc = FeatureCollection{}
	for key, value := range temp {
		switch key {
		case "type":
			fc.Type, _ = bson.RawValue(value).StringValueOK()
		case "bbox":
			if err = value.Unmarshal(&fc.BBox); err != nil {
				return
			}
		case "features":
			if err = value.Unmarshal(&fc.Features); err != nil {
				return
			}
		default:
			if fc.ExtraMembers == nil {
				fc.ExtraMembers = Properties{}
			}

			var val interface{}
			if err = value.Unmarshal(&val); err != nil {
				return
			} else {
				fc.ExtraMembers[key] = val
			}
		}
	}

	if fc.Type != featureCollection {
		return fmt.Errorf("geojson: not a feature collection: type=%s", fc.Type)
	}

	return nil
}

// UnmarshalFeatureCollection decodes the data into a GeoJSON feature collection.
// Alternately one can call json.Unmarshal(fc) directly for the same result.
func UnmarshalFeatureCollection(data []byte) (*FeatureCollection, error) {
	fc := &FeatureCollection{}
	if err := fc.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	return fc, nil
}

// Append appends a feature to the collection.
func (fc *FeatureCollection) Append(feature *Feature) *FeatureCollection {
	fc.Features = append(fc.Features, feature)
	return fc
}

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
