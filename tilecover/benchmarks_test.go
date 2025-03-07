package tilecover

import (
	"testing"

	"github.com/pchchv/geo"
)

func BenchmarkPoint(b *testing.B) {
	p := geo.Point{-77.15664982795715, 38.87419791355846}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(p, 6); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRoad_z6(b *testing.B) {
	g := loadFeature(b, "./testdata/road.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 6); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRoad_z18(b *testing.B) {
	g := loadFeature(b, "./testdata/road.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 18); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRoad_z28(b *testing.B) {
	g := loadFeature(b, "./testdata/road.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 28); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussia_z6(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 6); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussia_z8(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 8); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussia_z10(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 10); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussia_z0z9(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tiles, _ := Geometry(g, 9)
		MergeUp(tiles, 0)
	}
}

func BenchmarkRussiaLine_z6(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	g = geo.LineString(g.(geo.Polygon)[0])
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 6); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussiaLine_z8(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	g = geo.LineString(g.(geo.Polygon)[0])
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 8); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}

func BenchmarkRussiaLine_z10(b *testing.B) {
	g := loadFeature(b, "./testdata/russia.geojson").Geometry
	g = geo.LineString(g.(geo.Polygon)[0])
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Geometry(g, 10); err != nil {
			b.Fatalf("unexpected error: %e", err)
		}
	}
}
