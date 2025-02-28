package geojson_test

import (
	"fmt"
	"log"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/planar"
	"github.com/pchchv/geo/quadtree"
)

type CentroidPoint struct {
	*geojson.Feature
}

func (cp CentroidPoint) Point() geo.Point {
	// this is where you would decide how to define
	// the representative point of the feature.
	c, _ := planar.CentroidArea(cp.Geometry)
	return c
}

func Example_centroid() {
	qt := quadtree.New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	// feature with center {0.5, 0.5} but centroid {0.25, 0.25}
	f := geojson.NewFeature(geo.MultiPoint{{0, 0}, {0, 0}, {0, 0}, {1, 1}})
	f.Properties["centroid"] = "0.25"
	if err := qt.Add(CentroidPoint{f}); err != nil {
		log.Fatalf("unexpected error: %e", err)
	}

	// feature with centroid {0.6, 0.6}
	f = geojson.NewFeature(geo.Point{0.6, 0.6})
	f.Properties["centroid"] = "0.6"
	if err := qt.Add(CentroidPoint{f}); err != nil {
		log.Fatalf("unexpected error: %e", err)
	}

	feature := qt.Find(geo.Point{0.5, 0.5}).(CentroidPoint).Feature
	fmt.Printf("centroid=%s", feature.Properties["centroid"])

	// Output:
	// centroid=0.6
}
