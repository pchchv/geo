package geojson

import (
	"errors"

	"github.com/pchchv/geo"
	"go.mongodb.org/mongo-driver/bson"
)

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

// LineString is a helper type that
// will marshal to/from a GeoJSON LineString geometry.
type LineString geo.LineString

// MarshalJSON will convert the LineString into a GeoJSON LineString geometry.
func (ls LineString) MarshalJSON() ([]byte, error) {
	return marshalJSON(&Geometry{Coordinates: geo.LineString(ls)})
}

// MarshalBSON will convert the LineString into a GeoJSON LineString geometry.
func (ls LineString) MarshalBSON() ([]byte, error) {
	return bson.Marshal(&Geometry{Coordinates: geo.LineString(ls)})
}

// UnmarshalJSON will unmarshal the GeoJSON MultiPoint geometry.
func (ls *LineString) UnmarshalJSON(data []byte) error {
	g := &Geometry{}
	err := unmarshalJSON(data, &g)
	if err != nil {
		return err
	}

	lineString, ok := g.Coordinates.(geo.LineString)
	if !ok {
		return errors.New("geojson: not a LineString type")
	}

	*ls = LineString(lineString)
	return nil
}

// UnmarshalBSON will unmarshal the GeoJSON MultiPoint geometry.
func (ls *LineString) UnmarshalBSON(data []byte) error {
	g := &Geometry{}
	err := bson.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	lineString, ok := g.Coordinates.(geo.LineString)
	if !ok {
		return errors.New("geojson: not a LineString type")
	}

	*ls = LineString(lineString)
	return nil
}

// Geometry will return the geo.Geometry version of the data.
func (ls LineString) Geometry() geo.Geometry {
	return geo.LineString(ls)
}

// MultiLineString is a helper type that
// will marshal to/from a GeoJSON MultiLineString geometry.
type MultiLineString geo.MultiLineString

// MarshalJSON will convert the MultiLineString into a GeoJSON MultiLineString geometry.
func (mls MultiLineString) MarshalJSON() ([]byte, error) {
	return marshalJSON(&Geometry{Coordinates: geo.MultiLineString(mls)})
}

// MarshalBSON will convert the MultiLineString into a GeoJSON MultiLineString geometry.
func (mls MultiLineString) MarshalBSON() ([]byte, error) {
	return bson.Marshal(&Geometry{Coordinates: geo.MultiLineString(mls)})
}

// UnmarshalJSON will unmarshal the GeoJSON MultiPoint geometry.
func (mls *MultiLineString) UnmarshalJSON(data []byte) error {
	g := &Geometry{}
	err := unmarshalJSON(data, &g)
	if err != nil {
		return err
	}

	multilineString, ok := g.Coordinates.(geo.MultiLineString)
	if !ok {
		return errors.New("geojson: not a MultiLineString type")
	}

	*mls = MultiLineString(multilineString)
	return nil
}

// UnmarshalBSON will unmarshal the GeoJSON MultiPoint geometry.
func (mls *MultiLineString) UnmarshalBSON(data []byte) error {
	g := &Geometry{}
	err := bson.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	multilineString, ok := g.Coordinates.(geo.MultiLineString)
	if !ok {
		return errors.New("geojson: not a MultiLineString type")
	}

	*mls = MultiLineString(multilineString)
	return nil
}

// Geometry will return the geo.Geometry version of the data.
func (mls MultiLineString) Geometry() geo.Geometry {
	return geo.MultiLineString(mls)
}

// Polygon is a helper type that
// will marshal to/from a GeoJSON Polygon geometry.
type Polygon geo.Polygon

// MarshalJSON will convert the Polygon into a GeoJSON Polygon geometry.
func (p Polygon) MarshalJSON() ([]byte, error) {
	return marshalJSON(&Geometry{Coordinates: geo.Polygon(p)})
}

// MarshalBSON will convert the Polygon into a GeoJSON Polygon geometry.
func (p Polygon) MarshalBSON() ([]byte, error) {
	return bson.Marshal(&Geometry{Coordinates: geo.Polygon(p)})
}

// UnmarshalJSON will unmarshal the GeoJSON Polygon geometry.
func (p *Polygon) UnmarshalJSON(data []byte) error {
	g := &Geometry{}
	err := unmarshalJSON(data, &g)
	if err != nil {
		return err
	}

	polygon, ok := g.Coordinates.(geo.Polygon)
	if !ok {
		return errors.New("geojson: not a Polygon type")
	}

	*p = Polygon(polygon)
	return nil
}

// UnmarshalBSON will unmarshal the GeoJSON Polygon geometry.
func (p *Polygon) UnmarshalBSON(data []byte) error {
	g := &Geometry{}
	err := bson.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	polygon, ok := g.Coordinates.(geo.Polygon)
	if !ok {
		return errors.New("geojson: not a Polygon type")
	}

	*p = Polygon(polygon)
	return nil
}

// Geometry will return the geo.Geometry version of the data.
func (p Polygon) Geometry() geo.Geometry {
	return geo.Polygon(p)
}
