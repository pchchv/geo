package project

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestGeometry(t *testing.T) {
	for _, g := range geo.AllGeometries {
		// should not panic with unsupported type
		Geometry(g, Mercator.ToWGS84)
	}
}
