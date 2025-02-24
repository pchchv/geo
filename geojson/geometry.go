package geojson

import (
	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

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

// MarshalJSON will marshal the geometry into the correct JSON structure.
func (g *Geometry) MarshalJSON() ([]byte, error) {
	if g.Coordinates == nil && len(g.Geometries) == 0 {
		return []byte(`null`), nil
	}

	ng := newGeometryMarshallDoc(g)
	return marshalJSON(ng)
}

// MarshalBSON will convert the geometry into a
// BSON document with the structure of a GeoJSON Geometry.
// This function is used when the geometry is the
// top level document to be marshalled.
func (g *Geometry) MarshalBSON() ([]byte, error) {
	ng := newGeometryMarshallDoc(g)
	return bson.Marshal(ng)
}

// MarshalBSONValue will marshal the geometry into a
// BSON value with the structure of a GeoJSON Geometry.
func (g *Geometry) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// implementing MarshalBSONValue allows us to
	// marshal into a null value needed to
	// match behavior with the JSON marshalling
	if g.Coordinates == nil && len(g.Geometries) == 0 {
		return bsontype.Type(0x0A), nil, nil
	}

	ng := newGeometryMarshallDoc(g)
	return bson.MarshalValue(ng)
}

type geometryMarshallDoc struct {
	Type        string       `json:"type" bson:"type"`
	Coordinates geo.Geometry `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	Geometries  []*Geometry  `json:"geometries,omitempty" bson:"geometries,omitempty"`
}

type bsonGeometry struct {
	Type        string        `json:"type" bson:"type"`
	Coordinates bson.RawValue `json:"coordinates" bson:"coordinates"`
	Geometries  []*Geometry   `json:"geometries,omitempty" bson:"geometries"`
}

type jsonGeometry struct {
	Type        string           `json:"type"`
	Coordinates nocopyRawMessage `json:"coordinates"`
	Geometries  []*Geometry      `json:"geometries,omitempty"`
}

func newGeometryMarshallDoc(g *Geometry) *geometryMarshallDoc {
	ng := &geometryMarshallDoc{}
	switch g := g.Coordinates.(type) {
	case geo.Ring:
		ng.Coordinates = geo.Polygon{g}
	case geo.Bound:
		ng.Coordinates = g.ToPolygon()
	case geo.Collection:
		ng.Geometries = make([]*Geometry, 0, len(g))
		for _, c := range g {
			ng.Geometries = append(ng.Geometries, NewGeometry(c))
		}
		ng.Type = g.GeoJSONType()
	default:
		ng.Coordinates = g
	}

	if ng.Coordinates != nil {
		ng.Type = ng.Coordinates.GeoJSONType()
	}

	if len(g.Geometries) > 0 {
		ng.Geometries = g.Geometries
		ng.Type = geo.Collection{}.GeoJSONType()
	}

	return ng
}
