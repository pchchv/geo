package planar_test

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

func ExampleArea() {
	// +
	// |\
	// | \
	// |  \
	// +---+

	r := geo.Ring{{0, 0}, {3, 0}, {0, 4}, {0, 0}}
	a := planar.Area(r)

	fmt.Println(a)
	// Output:
	// 6
}

func ExampleDistance() {
	d := planar.Distance(geo.Point{0, 0}, geo.Point{3, 4})

	fmt.Println(d)
	// Output:
	// 5
}

func ExampleLength() {
	// +
	// |\
	// | \
	// |  \
	// +---+

	r := geo.Ring{{0, 0}, {3, 0}, {0, 4}, {0, 0}}
	l := planar.Length(r)

	fmt.Println(l)
	// Output:
	// 12
}
