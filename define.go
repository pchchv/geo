package geo

// Constants to define orientation.
// They follow the right hand rule for orientation.
const (
	CW  Orientation = -1 // stands for Clock Wise
	CCW Orientation = 1  // stands for Counter Clock Wise
)

// Pointer is something that can be represented by a point.
type Pointer interface {
	Point() Point
}

// Orientation defines the order of the points in a polygon or closed ring.
type Orientation int8

// Projection moves a point from one space to another.
type Projection func(Point) Point

// DistanceFunc computes the distance between two points.
type DistanceFunc func(Point, Point) float64

// Simplifier can simplify geometry.
type Simplifier interface {
	Simplify(g Geometry) Geometry
	LineString(ls LineString) LineString
	MultiLineString(mls MultiLineString) MultiLineString
	Ring(r Ring) Ring
	Polygon(p Polygon) Polygon
	MultiPolygon(mp MultiPolygon) MultiPolygon
	Collection(c Collection) Collection
}
