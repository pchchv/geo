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
