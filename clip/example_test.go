package clip_test

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/clip"
)

func ExampleGeometry() {
	bound := geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{30, 30}}
	ls := geo.LineString{
		{-10, 10}, {10, 10}, {10, -10}, {20, -10}, {20, 10},
		{40, 10}, {40, 20}, {20, 20}, {20, 40}, {10, 40},
		{10, 20}, {5, 20}, {-10, 20}}

	// returns an geo.Geometry interface
	clipped := clip.Geometry(bound, ls)

	fmt.Println(clipped)
	// Output:
	// [[[0 10] [10 10] [10 0]] [[20 0] [20 10] [30 10]] [[30 20] [20 20] [20 30]] [[10 30] [10 20] [5 20] [0 20]]]
}
