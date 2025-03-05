package simplifier_test

import (
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
	"github.com/pchchv/geo/simplifier"
)

func ExampleDouglasPeuckerSimplifier() {
	//  +
	//   \
	//    \
	//     +
	//      \
	//       \
	//  +-----+
	original := geo.LineString{{0, 0}, {2, 0}, {1, 1}, {0, 2}}

	// low threshold just removes the colinear point
	reduced := simplifier.DouglasPeucker(0.0).Simplify(original.Clone())
	fmt.Println(reduced)

	// high threshold just leaves start and end
	reduced = simplifier.DouglasPeucker(2).Simplify(original)
	fmt.Println(reduced)

	// Output:
	// [[0 0] [2 0] [0 2]]
	// [[0 0] [0 2]]
}

func ExampleRadialSimplifier() {
	//  +
	//   \
	//    \
	//     +
	//     |
	//  +--+
	original := geo.LineString{{0, 0}, {1, 0}, {1, 1}, {0, 2}}

	// will remove the points within 1.0 of the previous point
	// in this case just the second point
	reduced := simplifier.Radial(planar.Distance, 1.0).Simplify(original.Clone())
	fmt.Println(reduced)

	// will remove the 2nd and 3rd point since it's within 1.5 or the first point.
	reduced = simplifier.Radial(planar.Distance, 1.5).Simplify(original)
	fmt.Println(reduced)
}

func ExampleVisvalingamSimplifier() {
	original := geo.Ring{}
	threshold := 0.5 // define a threshold value

	// will remove all whose triangle is smaller than `threshold`
	reduced := simplifier.VisvalingamThreshold(threshold).Simplify(original)
	fmt.Println(reduced)

	toKeep := 3 // define the number of points to keep
	// will remove points until there are only `toKeep` points left.
	reduced = simplifier.VisvalingamKeep(toKeep).Simplify(original)
	fmt.Println(reduced)

	// One can also combine the parameters.
	// Will continue to remove points until:
	//   - there are no more below the threshold,
	//   - or the new path is of length `toKeep`
	reduced = simplifier.Visvalingam(threshold, toKeep).Simplify(original)
	fmt.Println(reduced)
}
