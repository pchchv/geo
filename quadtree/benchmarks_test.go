package quadtree

import (
	"math"
	"math/rand"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

func BenchmarkAdd(b *testing.B) {
	r := rand.New(rand.NewSource(22))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := qt.Add(geo.Point{r.Float64(), r.Float64()}); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkRandomFind1000(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		if err := qt.Add(geo.Point{r.Float64(), r.Float64()}); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Find(geo.Point{r.Float64(), r.Float64()})
	}
}

func BenchmarkRandomFind1000Naive(b *testing.B) {
	points := []geo.Point{}
	r := rand.New(rand.NewSource(42))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		if err := qt.Add(p); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		points = append(points, p)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var best geo.Point
		min := math.MaxFloat64
		looking := geo.Point{r.Float64(), r.Float64()}
		for _, p := range points {
			if d := planar.DistanceSquared(looking, p); d < min {
				min, best = d, p
			}
		}

		_ = best
	}
}

func BenchmarkRandomInBound1000(b *testing.B) {
	r := rand.New(rand.NewSource(43))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		if err := qt.Add(p); err != nil {
			b.Fatalf("unexpected error for %v: %v", p, err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		qt.InBound(nil, p.Bound().Pad(0.1))
	}
}

func BenchmarkRandomInBound1000Naive(b *testing.B) {
	points := []geo.Point{}
	r := rand.New(rand.NewSource(43))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		if err := qt.Add(p); err != nil {
			b.Fatalf("unexpected error for %v: %v", p, err)
		}

		points = append(points, p)
	}

	var near []geo.Point
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		b := geo.Bound{Min: p, Max: p}
		b = b.Pad(0.1)
		near = near[:0]
		for _, p := range points {
			if b.Contains(p) {
				near = append(near, p)
			}
		}

		_ = len(near)
	}
}

func BenchmarkRandomInBound1000Buf(b *testing.B) {
	r := rand.New(rand.NewSource(43))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		if err := qt.Add(p); err != nil {
			b.Fatalf("unexpected error for %v: %v", p, err)
		}
	}

	var buf []geo.Pointer
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		buf = qt.InBound(buf, p.Bound().Pad(0.1))
	}
}
