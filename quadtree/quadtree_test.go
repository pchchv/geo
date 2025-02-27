package quadtree

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestNew(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{0, 2}, Max: geo.Point{1, 3}}
	qt := New(bound)

	if !qt.Bound().Equal(bound) {
		t.Errorf("should use provided bound, got %v", qt.Bound())
	}
}
