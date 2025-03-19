package wkbcommon

import (
	"bytes"
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

func compare(t testing.TB, e geo.Geometry, s int, b []byte) {
	t.Helper()
	// Decoder
	g, srid, err := NewDecoder(bytes.NewReader(b)).Decode()
	if err != nil {
		t.Fatalf("decoder: read error: %e", err)
	}

	if !geo.Equal(g, e) {
		t.Errorf("decoder: incorrect geometry: %v != %v", g, e)
	}

	if srid != s {
		t.Errorf("decoder: incorrect srid: %v != %v", srid, s)
	}

	// Umarshal
	g, srid, err = Unmarshal(b)
	if err != nil {
		t.Fatalf("unmarshal: read error: %e", err)
	}

	if !geo.Equal(g, e) {
		t.Errorf("unmarshal: incorrect geometry: %v != %v", g, e)
	}

	if srid != s {
		t.Errorf("decoder: incorrect srid: %v != %v", srid, s)
	}

	var data []byte
	if b[0] == 0 {
		data, err = Marshal(g, s, binary.BigEndian)
	} else {
		data, err = Marshal(g, s, binary.LittleEndian)
	}
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	}

	if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("marshal: incorrent encoding")
	}

	// preallocation
	if l := GeomLength(e, srid != 0); len(data) != l {
		t.Errorf("prealloc length: %v != %v", len(data), l)
	}
}
