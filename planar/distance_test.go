package planar

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestDistance(t *testing.T) {
	p1, p2 := geo.Point{0, 0}, geo.Point{3, 4}
	if d := Distance(p1, p2); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}

	if d := Distance(p2, p1); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}
}
