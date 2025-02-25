package geojson

import (
	"errors"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var ErrInvalidGeometry = errors.New("geojson: invalid geometry")

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

// Geometry returns the geo.Geometry for the geojson Geometry.
// This will convert the "Geometries" into a geo.Collection if applicable.
func (g *Geometry) Geometry() geo.Geometry {
	if g.Coordinates != nil {
		return g.Coordinates
	}

	c := make(geo.Collection, 0, len(g.Geometries))
	for _, geom := range g.Geometries {
		c = append(c, geom.Geometry())
	}

	return c
}

// UnmarshalGeometry decodes the JSON data into a GeoJSON feature.
// Alternately one can call json.Unmarshal(g) directly for the same result.
func UnmarshalGeometry(data []byte) (g *Geometry, err error) {
	if err = unmarshalJSON(data, g); err != nil {
		return nil, err
	}

	return
}

// UnmarshalJSON will unmarshal the correct geometry from the JSON structure.
func (g *Geometry) UnmarshalJSON(data []byte) (err error) {
	jg := &jsonGeometry{}
	if err = unmarshalJSON(data, jg); err != nil {
		return
	}

	switch jg.Type {
	case "Point":
		p := geo.Point{}
		if err = unmarshalJSON(jg.Coordinates, &p); err != nil {
			return
		}
		g.Coordinates = p
	case "MultiPoint":
		mp := geo.MultiPoint{}
		if err = unmarshalJSON(jg.Coordinates, &mp); err != nil {
			return
		}
		g.Coordinates = mp
	case "LineString":
		ls := geo.LineString{}
		if err = unmarshalJSON(jg.Coordinates, &ls); err != nil {
			return
		}
		g.Coordinates = ls
	case "MultiLineString":
		mls := geo.MultiLineString{}
		if err = unmarshalJSON(jg.Coordinates, &mls); err != nil {
			return
		}
		g.Coordinates = mls
	case "Polygon":
		p := geo.Polygon{}
		if err = unmarshalJSON(jg.Coordinates, &p); err != nil {
			return
		}
		g.Coordinates = p
	case "MultiPolygon":
		mp := geo.MultiPolygon{}
		if err = unmarshalJSON(jg.Coordinates, &mp); err != nil {
			return
		}
		g.Coordinates = mp
	case "GeometryCollection":
		g.Geometries = jg.Geometries
	default:
		return ErrInvalidGeometry
	}

	g.Type = g.Geometry().GeoJSONType()
	return nil
}

// UnmarshalBSON will unmarshal a BSON document created with bson.Marshal.
func (g *Geometry) UnmarshalBSON(data []byte) (err error) {
	bg := &bsonGeometry{}
	if err = bson.Unmarshal(data, bg); err != nil {
		return
	}

	switch bg.Type {
	case "Point":
		p := geo.Point{}
		if err = bg.Coordinates.Unmarshal(&p); err != nil {
			return
		}
		g.Coordinates = p
	case "MultiPoint":
		mp := geo.MultiPoint{}
		if err = bg.Coordinates.Unmarshal(&mp); err != nil {
			return
		}
		g.Coordinates = mp
	case "LineString":
		ls := geo.LineString{}
		if err = bg.Coordinates.Unmarshal(&ls); err != nil {
			return
		}
		g.Coordinates = ls
	case "MultiLineString":
		mls := geo.MultiLineString{}
		if err = bg.Coordinates.Unmarshal(&mls); err != nil {
			return
		}
		g.Coordinates = mls
	case "Polygon":
		p := geo.Polygon{}
		if err = bg.Coordinates.Unmarshal(&p); err != nil {
			return
		}
		g.Coordinates = p
	case "MultiPolygon":
		mp := geo.MultiPolygon{}
		if err = bg.Coordinates.Unmarshal(&mp); err != nil {
			return
		}
		g.Coordinates = mp
	case "GeometryCollection":
		g.Geometries = bg.Geometries
	default:
		return ErrInvalidGeometry
	}

	g.Type = g.Geometry().GeoJSONType()
	return nil
}

// Point is a helper type that
// will marshal to/from a GeoJSON Point geometry.
type Point geo.Point

// MarshalJSON will convert the Point into a GeoJSON Point geometry.
func (p Point) MarshalJSON() ([]byte, error) {
	return marshalJSON(&Geometry{Coordinates: geo.Point(p)})
}

// MarshalBSON will convert the Point into a
// BSON value following the GeoJSON Point structure.
func (p Point) MarshalBSON() ([]byte, error) {
	return bson.Marshal(&Geometry{Coordinates: geo.Point(p)})
}

// UnmarshalJSON will unmarshal the GeoJSON Point geometry.
func (p *Point) UnmarshalJSON(data []byte) error {
	g := &Geometry{}
	err := unmarshalJSON(data, &g)
	if err != nil {
		return err
	}

	point, ok := g.Coordinates.(geo.Point)
	if !ok {
		return errors.New("geojson: not a Point type")
	}

	*p = Point(point)
	return nil
}

// UnmarshalBSON will unmarshal GeoJSON Point geometry.
func (p *Point) UnmarshalBSON(data []byte) error {
	g := &Geometry{}
	err := bson.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	point, ok := g.Coordinates.(geo.Point)
	if !ok {
		return errors.New("geojson: not a Point type")
	}

	*p = Point(point)
	return nil
}

// Geometry will return the geo.Geometry version of the data.
func (p Point) Geometry() geo.Geometry {
	return geo.Point(p)
}

// MultiPoint is a helper type that
// will marshal to/from a GeoJSON MultiPoint geometry.
type MultiPoint geo.MultiPoint

// MarshalJSON will convert the MultiPoint into a GeoJSON MultiPoint geometry.
func (mp MultiPoint) MarshalJSON() ([]byte, error) {
	return marshalJSON(&Geometry{Coordinates: geo.MultiPoint(mp)})
}

// MarshalBSON will convert the MultiPoint into a GeoJSON MultiPoint geometry BSON.
func (mp MultiPoint) MarshalBSON() ([]byte, error) {
	return bson.Marshal(&Geometry{Coordinates: geo.MultiPoint(mp)})
}

// UnmarshalJSON will unmarshal the GeoJSON MultiPoint geometry.
func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	g := &Geometry{}
	err := unmarshalJSON(data, &g)
	if err != nil {
		return err
	}

	multiPoint, ok := g.Coordinates.(geo.MultiPoint)
	if !ok {
		return errors.New("geojson: not a MultiPoint type")
	}

	*mp = MultiPoint(multiPoint)
	return nil
}

// UnmarshalBSON will unmarshal the GeoJSON MultiPoint geometry.
func (mp *MultiPoint) UnmarshalBSON(data []byte) error {
	g := &Geometry{}
	err := bson.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	multiPoint, ok := g.Coordinates.(geo.MultiPoint)
	if !ok {
		return errors.New("geojson: not a MultiPoint type")
	}

	*mp = MultiPoint(multiPoint)
	return nil
}

// Geometry will return the geo.Geometry version of the data.
func (mp MultiPoint) Geometry() geo.Geometry {
	return geo.MultiPoint(mp)
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
