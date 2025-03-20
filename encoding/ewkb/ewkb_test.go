package ewkb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
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

func BenchmarkEncode_Point(b *testing.B) {
	g := geo.Point{1, 2}
	e := NewEncoder(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := e.Encode(g)
		if err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkEncode_LineString(b *testing.B) {
	g := geo.LineString{
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
	}
	e := NewEncoder(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := e.Encode(g)
		if err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func MustDecodeHex(s string) []byte {
	if b, err := hex.DecodeString(s); err != nil {
		panic(err)
	} else {
		return b
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

	// Marshal
	var data []byte
	if b[0] == 0 {
		data, err = Marshal(g, srid, binary.BigEndian)
	} else {
		data, err = Marshal(g, srid, binary.LittleEndian)
	}
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	}

	if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("marshal: incorrent encoding")
	}

	// Encode
	buf := bytes.NewBuffer(nil)
	en := NewEncoder(buf)
	if b[0] == 0 {
		en.SetByteOrder(binary.BigEndian)
	} else {
		en.SetByteOrder(binary.LittleEndian)
	}

	en.SetSRID(s)
	if err = en.Encode(e); err != nil {
		t.Errorf("encode error: %e", err)
	}

	if !bytes.Equal(data, buf.Bytes()) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("encode: incorrent encoding")
	}

	// pass in srid
	buf.Reset()
	en.SetSRID(10101)
	if err = en.Encode(e, s); err != nil {
		t.Errorf("encode with srid error: %e", err)
	}

	if !bytes.Equal(data, buf.Bytes()) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("encode with srid: incorrent encoding")
	}

	// preallocation
	if l := wkbcommon.GeomLength(e, s != 0); len(data) != l {
		t.Errorf("prealloc length: %v != %v", len(data), l)
	}
}
