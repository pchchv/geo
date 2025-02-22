package geo

// compile time checks
var (
	_ Geometry = Ring{}
	_ Geometry = Point{}
	_ Geometry = Bound{}
	_ Geometry = Polygon{}
	_ Geometry = Collection{}
	_ Geometry = MultiPoint{}
	_ Geometry = LineString{}
	_ Geometry = MultiPolygon{}
	_ Geometry = MultiLineString{}
)

// Geometry represents the shared attributes of a geometry.
type Geometry interface {
	GeoJSONType() string
	Dimensions() int // i.e., 0d, 1d, 2d
	Bound() Bound
	private() // requiring because sub package type switch over all possible types
}

func (r Ring) private()              {}
func (p Point) private()             {}
func (b Bound) private()             {}
func (p Polygon) private()           {}
func (c Collection) private()        {}
func (mp MultiPoint) private()       {}
func (ls LineString) private()       {}
func (mp MultiPolygon) private()     {}
func (mls MultiLineString) private() {}

// Collection is a collection of geometries that is also a Geometry.
type Collection []Geometry

// Bound returns the bounding box of all the Geometries combined.
func (c Collection) Bound() Bound {
	if len(c) == 0 {
		return emptyBound
	}

	var b Bound
	start := -1
	for i, g := range c {
		if g != nil {
			start = i
			b = g.Bound()
			break
		}
	}

	if start == -1 {
		return emptyBound
	}

	for i := start + 1; i < len(c); i++ {
		if c[i] == nil {
			continue
		}

		b = b.Union(c[i].Bound())
	}

	return b
}

// Dimensions returns the max of the dimensions of the collection.
func (c Collection) Dimensions() (max int) {
	max--
	for _, g := range c {
		if d := g.Dimensions(); d > max {
			max = d
		}
	}

	return
}

// GeoJSONType returns the geometry collection type.
func (c Collection) GeoJSONType() string {
	return "GeometryCollection"
}
