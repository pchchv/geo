package geometries

import (
	"math"
	"testing"

	"github.com/pchchv/geo"
)

const epsilon = 1e-6

func TestDistance(t *testing.T) {
	p1 := geo.Point{-1.8444, 53.1506}
	p2 := geo.Point{0.1406, 52.2047}
	if d := Distance(p1, p2); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = geo.Point{0.5, 30}
	p2 = geo.Point{-0.5, 30}
	dFast := Distance(p1, p2)
	p1 = geo.Point{179.5, 30}
	p2 = geo.Point{-179.5, 30}
	if d := Distance(p1, p2); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestDistanceHaversine(t *testing.T) {
	p1 := geo.Point{-1.8444, 53.1506}
	p2 := geo.Point{0.1406, 52.2047}
	if d := DistanceHaversine(p1, p2); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = geo.Point{0.5, 30}
	p2 = geo.Point{-0.5, 30}
	dHav := DistanceHaversine(p1, p2)
	p1 = geo.Point{179.5, 30}
	p2 = geo.Point{-179.5, 30}
	if d := DistanceHaversine(p1, p2); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}
