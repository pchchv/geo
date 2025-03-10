package geo

import (
	"encoding/json"
	"testing"
)

var (
	pointData           = []byte(`[-81.60540868,41.51522539]`)
	polygonShortData    = append(append([]byte{'['}, lineStringShortData...), byte(']'))
	polygonLongData     = append(append([]byte{'['}, lineStringLongData...), byte(']'))
	lineStringShortData = []byte(`[
		[-81.60540868,41.51522539],[-81.6049915,41.51523057],
		[-81.60499258,41.51528156],[-81.60506867,41.51528061],
		[-81.60540868,41.51522539]
	]`)
	lineStringLongData = []byte(`[
		[-81.60540868,41.51522539],[-81.6049915,41.51523057],
		[-81.60499258,41.51528156],[-81.60506867,41.51528061],
		[-81.60507145,41.5154027],[-81.60500219,41.51540351],
		[-81.60500345,41.51545786],[-81.60529819,41.51545416],
		[-81.60530133,41.51559117],[-81.60487193,41.51559662],
		[-81.60486062,41.51508837],[-81.60521401,41.515084],
		[-81.60521123,41.51495707],[-81.60480788,41.51496219],
		[-81.6047998,41.51460138],[-81.60496455,41.51459936],
		[-81.60496778,41.51474909],[-81.6053979,41.51474371],
		[-81.60540868,41.51522539]
	]`)
)

func TestPointJSON(t *testing.T) {
	p1 := Point{1, 2.1}
	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine: %e", err)
	}

	if string(data) != "[1,2.1]" {
		t.Errorf("incorrect json: %v", string(data))
	}

	var p2 Point
	if err = json.Unmarshal(data, &p2); err != nil {
		t.Errorf("should unmarshal just fine: %e", err)
	}

	if !p1.Equal(p2) {
		t.Errorf("not equal: %v", p2)
	}
}

func TestLineStringJSON(t *testing.T) {
	ls1 := LineString{{1.5, 2.5}, {3.5, 4.5}, {5.5, 6.5}}
	data, err := json.Marshal(ls1)
	if err != nil {
		t.Fatalf("should marshal just fine: %e", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5],[5.5,6.5]]" {
		t.Errorf("incorrect data: %v", string(data))
	}

	var ls2 LineString
	if err = json.Unmarshal(data, &ls2); err != nil {
		t.Fatalf("should unmarshal just fine: %e", err)
	}

	if !ls1.Equal(ls2) {
		t.Errorf("unmarshal not equal: %v", ls2)
	}

	// empty line
	ls1 = LineString{}
	if data, err = json.Marshal(ls1); err != nil {
		t.Errorf("should marshal just fine: %e", err)
	} else if string(data) != "[]" {
		t.Errorf("incorrect json: %v", string(data))
	}
}

func BenchmarkPointMarshalJSON(b *testing.B) {
	p := Point{}
	if err := json.Unmarshal(pointData, &p); err != nil {
		b.Fatalf("unable to unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkPointUnmarshalJSON(b *testing.B) {
	p := Point{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(pointData, &p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkLineStringMarshalJSON_short(b *testing.B) {
	ls := LineString{}
	if err := json.Unmarshal(lineStringShortData, &ls); err != nil {
		b.Fatalf("unable to unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(ls); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkLineStringUnmarshalJSON_short(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls := LineString{}
		if err := json.Unmarshal(lineStringShortData, &ls); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkLineStringMarshalJSON_long(b *testing.B) {
	ls := LineString{}
	if err := json.Unmarshal(lineStringLongData, &ls); err != nil {
		b.Fatalf("unable to unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(ls); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkLineStringUnmarshalJSON_long(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls := LineString{}
		if err := json.Unmarshal(lineStringLongData, &ls); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkPolygonMarshalJSON_short(b *testing.B) {
	p := Polygon{}
	if err := json.Unmarshal(polygonShortData, &p); err != nil {
		b.Fatalf("unable to unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkPolygonUnmarshalJSON_short(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := Polygon{}
		if err := json.Unmarshal(polygonShortData, &p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkPolygonMarshalJSON_long(b *testing.B) {
	p := Polygon{}
	if err := json.Unmarshal(polygonLongData, &p); err != nil {
		b.Fatalf("unable to unmarshal: %e", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}

func BenchmarkPolygonUnmarshalJSON_long(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := Polygon{}
		if err := json.Unmarshal(polygonLongData, &p); err != nil {
			b.Fatalf("marshal error: %e", err)
		}
	}
}
