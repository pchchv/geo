package wkbcommon

import (
	"encoding/binary"
	"testing"

	"github.com/pchchv/geo"
)

func TestMarshal(t *testing.T) {
	for _, g := range geo.AllGeometries {
		if _, err := Marshal(g, 0, binary.BigEndian); err != nil {
			t.Fatalf("unexpected error: %e", err)
		}
	}
}

func TestMustMarshal(t *testing.T) {
	for _, g := range geo.AllGeometries {
		MustMarshal(g, 0, binary.BigEndian)
	}
}
