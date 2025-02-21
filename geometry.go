package geo

// compile time checks
var (
	_ Geometry = Point{}
	_ Geometry = Bound{}
	_ Geometry = Collection{}
)

// Geometry represents the shared attributes of a geometry.
type Geometry interface {
	GeoJSONType() string
	Dimensions() int // i.e., 0d, 1d, 2d
	Bound() Bound
	private() // requiring because sub package type switch over all possible types
}

func (p Point) private()      {}
func (b Bound) private()      {}
func (c Collection) private() {}
