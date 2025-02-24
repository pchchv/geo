package geojson

import "github.com/pchchv/geo"

// Geometry matches the structure of a GeoJSON Geometry.
type Geometry struct {
	Type        string       `json:"type"`
	Coordinates geo.Geometry `json:"coordinates,omitempty"`
	Geometries  []*Geometry  `json:"geometries,omitempty"`
}

// NewGeometry will create a Geometry object but
// will convert the input into a GoeJSON geometry.
// ie. it will convert Rings and Bounds into Polygons.
func NewGeometry(g geo.Geometry) *Geometry {
	jg := &Geometry{}
	switch g := g.(type) {
	case geo.Ring:
		jg.Coordinates = geo.Polygon{g}
	case geo.Bound:
		jg.Coordinates = g.ToPolygon()
	case geo.Collection:
		for _, c := range g {
			jg.Geometries = append(jg.Geometries, NewGeometry(c))
		}
		jg.Type = g.GeoJSONType()
	default:
		jg.Coordinates = g
	}

	if jg.Coordinates != nil {
		jg.Type = jg.Coordinates.GeoJSONType()
	}

	return jg
}
