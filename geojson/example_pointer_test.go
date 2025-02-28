package geojson_test

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/planar"
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
