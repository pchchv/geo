package quadtree

import (
	"math/rand"
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

func TestNew(t *testing.T) {
	bound := geo.Bound{Min: geo.Point{0, 2}, Max: geo.Point{1, 3}}
	qt := New(bound)

	if !qt.Bound().Equal(bound) {
		t.Errorf("should use provided bound, got %v", qt.Bound())
	}
}

func TestQuadtreeFind(t *testing.T) {
	dim := 17
	points := geo.MultiPoint{}
	for i := 0; i < dim*dim; i++ {
		points = append(points, geo.Point{float64(i % dim), float64(i / dim)})
	}

	qt := New(points.Bound())
	for _, p := range points {
		err := qt.Add(p)
		if err != nil {
			t.Fatalf("unexpected error for %v: %v", p, err)
		}
	}

	cases := []struct {
		point    geo.Point
		expected geo.Point
	}{
		{point: geo.Point{0.1, 0.1}, expected: geo.Point{0, 0}},
		{point: geo.Point{3.1, 2.9}, expected: geo.Point{3, 3}},
		{point: geo.Point{7.1, 7.1}, expected: geo.Point{7, 7}},
		{point: geo.Point{0.1, 15.9}, expected: geo.Point{0, 16}},
		{point: geo.Point{15.9, 15.9}, expected: geo.Point{16, 16}},
	}

	for i, tc := range cases {
		if v := qt.Find(tc.point); !v.Point().Equal(tc.expected) {
			t.Errorf("incorrect point on %d, got %v", i, v)
		}
	}
}

func TestQuadtreeFind_Random(t *testing.T) {
	mp := geo.MultiPoint{}
	r := rand.New(rand.NewSource(42))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		mp = append(mp, geo.Point{r.Float64(), r.Float64()})
		if err := qt.Add(mp[i]); err != nil {
			t.Fatalf("unexpected error for %v: %v", mp[i], err)
		}
	}

	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		f := qt.Find(p)
		_, j := planar.DistanceFromWithIndex(mp, p)
		if e := mp[j]; !e.Equal(f.Point()) {
			t.Errorf("index: %d, unexpected point %v != %v", i, e, f.Point())
		}
	}
}
