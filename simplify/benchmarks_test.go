package simplify

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pchchv/geo"
)

func TestDouglasPeucker_BenchmarkData(t *testing.T) {
	cases := []struct {
		threshold float64
		length    int
	}{
		{0.1, 1118},
		{0.5, 257},
		{1.0, 144},
		{1.5, 95},
		{2.0, 71},
		{3.0, 46},
		{4.0, 39},
		{5.0, 33},
	}

	ls := benchmarkData()
	for i, tc := range cases {
		r := DouglasPeucker(tc.threshold).LineString(ls.Clone())
		if len(r) != tc.length {
			t.Errorf("%d: reduced poorly, %d != %d", i, len(r), tc.length)
		}
	}
}

func BenchmarkDouglasPeucker(b *testing.B) {
	var data []geo.LineString
	ls := benchmarkData()
	for i := 0; i < b.N; i++ {
		data = append(data, ls.Clone())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DouglasPeucker(0.1).LineString(data[i])
	}
}

func TestVisvalingam_BenchmarkData(t *testing.T) {
	cases := []struct {
		threshold float64
		length    int
	}{
		{0.1, 867},
		{0.5, 410},
		{1.0, 293},
		{1.5, 245},
		{2.0, 208},
		{3.0, 169},
		{4.0, 151},
		{5.0, 135},
	}

	ls := benchmarkData()
	for i, tc := range cases {
		r := VisvalingamThreshold(tc.threshold).LineString(ls.Clone())
		if len(r) != tc.length {
			t.Errorf("%d: data reduced poorly: %v != %v", i, len(r), tc.length)
		}
	}
}

func BenchmarkVisvalingam_Threshold(b *testing.B) {
	var data []geo.LineString
	ls := benchmarkData()
	for i := 0; i < b.N; i++ {
		data = append(data, ls.Clone())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VisvalingamThreshold(0.1).LineString(data[i])
	}
}

func BenchmarkVisvalingam_Keep(b *testing.B) {
	var data []geo.LineString
	ls := benchmarkData()
	toKeep := int(float64(len(ls)) / 1.616)
	for i := 0; i < b.N; i++ {
		data = append(data, ls.Clone())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VisvalingamKeep(toKeep).LineString(data[i])
	}
}

func benchmarkData() (ls geo.LineString) {
	f, err := os.Open("testdata/lisbon2portugal.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var points []float64
	if err = json.NewDecoder(f).Decode(&points); err != nil {
		panic(err)
	}

	for i := 0; i < len(points); i += 2 {
		ls = append(ls, geo.Point{points[i], points[i+1]})
	}

	return
}
