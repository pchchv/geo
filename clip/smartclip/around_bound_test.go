package smartclip

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestNexts(t *testing.T) {
	for i, next := range nexts[geo.CW] {
		if next != -1 && i != nexts[geo.CCW][next] {
			t.Errorf("incorrect %d: %d != %d", i, i, nexts[geo.CCW][next])
		}
	}
}
