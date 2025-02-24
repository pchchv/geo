package geojson

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestGeometry(t *testing.T) {
	for _, g := range geo.AllGeometries {
		NewGeometry(g)
	}
}
