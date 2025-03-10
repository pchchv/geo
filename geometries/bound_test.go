package geometries

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestBoundAroundPoint(t *testing.T) {
	p := geo.Point{
		5.42553,
		50.0359,
	}

	b := NewBoundAroundPoint(p, 1000000)
	if b.Center()[1] != p[1] {
		t.Errorf("should have correct center lat point")
	}

	if b.Center()[0] != p[0] {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if !b.Contains(geo.Point{3.412, 58.3838}) {
		t.Errorf("should have point included in bound")
	}

	b = NewBoundAroundPoint(p, 10000.0)
	if b.Center()[1] != p[1] {
		t.Errorf("should have correct center lat point")
	}

	if b.Center()[0] != p[0] {
		t.Errorf("should have correct center lon point")
	}

	//Given point is 968.9 km away from center
	if b.Contains(geo.Point{3.412, 58.3838}) {
		t.Errorf("should not have point included in bound")
	}
}
