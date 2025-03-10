package geometries_test

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geometries"
)

func ExampleArea() {
	poly := geo.Polygon{
		{
			{-122.4163816, 37.7792782},
			{-122.4162786, 37.7787626},
			{-122.4151027, 37.7789118},
			{-122.4152143, 37.7794274},
			{-122.4163816, 37.7792782},
		},
	}
	a := geometries.Area(poly)

	fmt.Printf("%f m^2", a)
	// Output:
	// 6073.368008 m^2
}
