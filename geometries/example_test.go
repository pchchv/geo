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

func ExampleDistance() {
	oakland := geo.Point{-122.270833, 37.804444}
	sf := geo.Point{-122.416667, 37.783333}
	d := geometries.Distance(oakland, sf)

	fmt.Printf("%0.3f meters", d)
	// Output:
	// 13042.047 meters
}

func ExampleLength() {
	poly := geo.Polygon{
		{
			{-122.4163816, 37.7792782},
			{-122.4162786, 37.7787626},
			{-122.4151027, 37.7789118},
			{-122.4152143, 37.7794274},
			{-122.4163816, 37.7792782},
		},
	}
	l := geometries.Length(poly)

	fmt.Printf("%0.0f meters", l)
	// Output:
	// 325 meters
}
