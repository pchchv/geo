package tilecover

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestRing_error(t *testing.T) {
	// not a closed ring
	f := loadFeature(t, "./testdata/line.geojson")
	l := f.Geometry.(geo.LineString)
	if _, err := Ring(geo.Ring(l), 25); err.Error() != "tilecover: uneven intersections, ring not closed?" {
		t.Errorf("incorrect error: %v", err)
	}
}
