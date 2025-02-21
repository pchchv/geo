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
