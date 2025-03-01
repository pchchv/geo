package quadtree_test

import (
	"fmt"
	"math/rand"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/quadtree"
)

func ExampleQuadtree_Find() {
	r := rand.New(rand.NewSource(42)) // to make things reproducible
	qt := quadtree.New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	// add 1000 random points
	for i := 0; i < 1000; i++ {
		if err := qt.Add(geo.Point{r.Float64(), r.Float64()}); err != nil {
			panic(err)
		}
	}

	nearest := qt.Find(geo.Point{0.5, 0.5})
	fmt.Printf("nearest: %+v\n", nearest)

	// Output:
	// nearest: [0.4930591659434973 0.5196585530161364]
}

func ExampleQuadtree_Matching() {
	r := rand.New(rand.NewSource(42)) // to make things reproducible
	type dataPoint struct {
		geo.Pointer
		visible bool
	}

	qt := quadtree.New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	// add 100 random points
	for i := 0; i < 100; i++ {
		if err := qt.Add(dataPoint{geo.Point{r.Float64(), r.Float64()}, false}); err != nil {
			panic(err)
		}
	}

	if err := qt.Add(dataPoint{geo.Point{0, 0}, true}); err != nil {
		panic(err)
	}

	nearest := qt.Matching(
		geo.Point{0.5, 0.5},
		func(p geo.Pointer) bool { return p.(dataPoint).visible },
	)

	fmt.Printf("nearest: %+v\n", nearest)

	// Output:
	// nearest: {Pointer:[0 0] visible:true}
}

func ExampleQuadtree_InBound() {
	r := rand.New(rand.NewSource(52)) // to make things reproducible
	qt := quadtree.New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	// add 1000 random points
	for i := 0; i < 1000; i++ {
		if err := qt.Add(geo.Point{r.Float64(), r.Float64()}); err != nil {
			panic(err)
		}
	}

	bounded := qt.InBound(nil, geo.Point{0.5, 0.5}.Bound().Pad(0.05))
	fmt.Printf("in bound: %v\n", len(bounded))

	// Output:
	// in bound: 10
}
