package resample

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestLineStringResampleEdgeCases(t *testing.T) {
	ls := geo.LineString{{0, 0}}
	_, ret := resampleEdgeCases(ls, 10)
	if !ret {
		t.Errorf("should return true")
	}

	// duplicate points
	ls = append(ls, geo.Point{0, 0})
	if ls, ret = resampleEdgeCases(ls, 10); !ret {
		t.Errorf("should return true")
	}

	if l := len(ls); l != 10 {
		t.Errorf("should reset to suggested points: %v != 10", l)
	}

	ls, _ = resampleEdgeCases(ls, 5)
	if l := len(ls); l != 5 {
		t.Errorf("should shorten if necessary: %v != 5", l)
	}
}
