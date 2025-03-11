package mvt

import (
	"math"
	"testing"

	"github.com/pchchv/geo/maptile"
)

const epsilon = 1e-6

func TestPowerOfTwoProjection(t *testing.T) {
	// verify tile coord does not overflow int32
	tile := maptile.New(1730576, 798477, 21)
	proj := newProjection(tile, 4096)
	center := tile.Center()
	planar := proj.ToTile(center)
	if planar[0] != 2048 {
		t.Errorf("incorrect lon projection: %v", planar[0])
	}

	if planar[1] != 2048 {
		t.Errorf("incorrect lat projection: %v", planar[1])
	}

	geo := proj.ToWGS84(planar)
	if math.Abs(geo[0]-center[0]) > epsilon {
		t.Errorf("lon miss match: %f != %f", geo[0], center[0])
	}

	if math.Abs(geo[1]-center[1]) > epsilon {
		t.Errorf("lat miss match: %f != %f", geo[1], center[1])
	}
}
