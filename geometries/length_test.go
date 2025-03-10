package geometries

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestLength(t *testing.T) {
	for _, g := range geo.AllGeometries {
		// should not panic with unsupported type
		Length(g)
	}
}

func TestLengthHaversine(t *testing.T) {
	for _, g := range geo.AllGeometries {
		// should not panic with unsupported type
		LengthHaversine(g)
	}
}
