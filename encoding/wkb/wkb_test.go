package wkb

import (
	"io"
	"testing"

	"github.com/pchchv/geo"
)

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
