package tilecover

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestGeometry(t *testing.T) {
	for _, g := range geo.AllGeometries {
		if _, err := Geometry(g, 1); err != nil {
			t.Fatalf("unexpected error for %T: %v", g, err)
		}
	}
}
