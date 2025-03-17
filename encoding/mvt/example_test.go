package mvt_test

import (
	"log"

	"github.com/pchchv/geo/encoding/mvt"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/geo/maptile"
	"github.com/pchchv/geo/simplifier"
)

func ExampleMarshal() {
	// Start with a set of feature collections defining each layer in lon/lat (WGS84).
	collections := map[string]*geojson.FeatureCollection{}
	// Convert to a layers object and project to tile coordinates.
	layers := mvt.NewLayers(collections)
	layers.ProjectToTile(maptile.New(17896, 24449, 16)) // x, y, z
	// Simplify the geometry now that it's in the tile coordinate space.
	layers.Simplify(simplifier.DouglasPeucker(1.0))
	// Depending on use-case remove empty geometry,
	// those two small to be represented in this tile space.
	// In this case lines shorter than 1, and areas smaller than 1.
	layers.RemoveEmpty(1.0, 1.0)
	// encoding using the Mapbox Vector Tile protobuf encoding.
	data, err := mvt.Marshal(layers) // this data is NOT gzipped.
	_ = data

	// error checking
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}

	// Sometimes MVT data is stored and transferred gzip compressed. In that case:
	data, err = mvt.MarshalGzipped(layers)
	_ = data

	// error checking
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}
}
