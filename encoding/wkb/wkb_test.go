package wkb

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

func TestMarshal(t *testing.T) {
	for _, g := range geo.AllGeometries {
		if _, err := Marshal(g, binary.BigEndian); err != nil {
			t.Fatalf("unexpected error: %e", err)
		}
	}
}

func TestMustMarshal(t *testing.T) {
	for _, g := range geo.AllGeometries {
		MustMarshal(g, binary.BigEndian)
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
		if err := e.Encode(g); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func compare(t testing.TB, e geo.Geometry, b []byte) {
	t.Helper()
	// Decoder
	g, err := NewDecoder(bytes.NewReader(b)).Decode()
	if err != nil {
		t.Fatalf("decoder: read error: %e", err)
	} else if !geo.Equal(g, e) {
		t.Errorf("decoder: incorrect geometry: %v != %v", g, e)
	}

	// Umarshal
	g, err = Unmarshal(b)
	if err != nil {
		t.Fatalf("unmarshal: read error: %e", err)
	} else if !geo.Equal(g, e) {
		t.Errorf("unmarshal: incorrect geometry: %v != %v", g, e)
	}

	var data []byte
	if b[0] == 0 {
		data, err = Marshal(g, binary.BigEndian)
	} else {
		data, err = Marshal(g, binary.LittleEndian)
	}
	if err != nil {
		t.Fatalf("marshal error: %e", err)
	} else if !bytes.Equal(data, b) {
		t.Logf("%v", data)
		t.Logf("%v", b)
		t.Errorf("marshal: incorrent encoding")
	}

	// preallocation
	if l := wkbcommon.GeomLength(e, false); len(data) != l {
		t.Errorf("prealloc length: %v != %v", len(data), l)
	}

	// Scanner
	var sg geo.Geometry
	switch e.(type) {
	case geo.Point:
		var p geo.Point
		err = Scanner(&p).Scan(b)
		sg = p
	case geo.MultiPoint:
		var mp geo.MultiPoint
		err = Scanner(&mp).Scan(b)
		sg = mp
	case geo.LineString:
		var ls geo.LineString
		err = Scanner(&ls).Scan(b)
		sg = ls
	case geo.MultiLineString:
		var mls geo.MultiLineString
		err = Scanner(&mls).Scan(b)
		sg = mls
	case geo.Polygon:
		var p geo.Polygon
		err = Scanner(&p).Scan(b)
		sg = p
	case geo.MultiPolygon:
		var mp geo.MultiPolygon
		err = Scanner(&mp).Scan(b)
		sg = mp
	case geo.Collection:
		var c geo.Collection
		err = Scanner(&c).Scan(b)
		sg = c
	default:
		t.Fatalf("unknown type: %T", e)
	}

	if err != nil {
		t.Errorf("scan error: %e", err)
	}

	if sg.GeoJSONType() != e.GeoJSONType() {
		t.Errorf("scanning to wrong type: %v != %v", sg.GeoJSONType(), e.GeoJSONType())
	}

	if !geo.Equal(sg, e) {
		t.Errorf("scan: incorrect geometry: %v != %v", sg, e)
	}
}
