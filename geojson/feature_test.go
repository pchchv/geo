package geojson

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestNewFeature(t *testing.T) {
	f := NewFeature(geo.Point{1, 2})
	if f.Type != "Feature" {
		t.Errorf("incorrect feature: %v != Feature", f.Type)
	}
}
