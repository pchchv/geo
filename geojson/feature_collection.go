package geojson

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
