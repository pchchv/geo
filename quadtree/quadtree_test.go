package quadtree

import (
	"math/rand"
	"reflect"
	"sort"
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

func TestQuadtreeMatching(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	if err := qt.Add(dataPointer{geo.Point{0, 0}, false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := qt.Add(dataPointer{geo.Point{1, 1}, true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		name     string
		filter   FilterFunc
		point    geo.Point
		expected geo.Pointer
	}{
		{
			name:     "no filtred",
			point:    geo.Point{0.1, 0.1},
			expected: geo.Point{0, 0},
		},
		{
			name:     "with filter",
			filter:   func(p geo.Pointer) bool { return p.(dataPointer).visible },
			point:    geo.Point{0.1, 0.1},
			expected: geo.Point{1, 1},
		},
		{
			name:     "match none filter",
			filter:   func(p geo.Pointer) bool { return false },
			point:    geo.Point{0.1, 0.1},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := qt.Matching(tc.point, tc.filter)
			// case 1: exact match, important for testing `nil`
			if v == tc.expected {
				return
			}

			// case 2: match on returned geo.Point value
			if !v.Point().Equal(tc.expected.Point()) {
				t.Errorf("incorrect point %v != %v", v, tc.expected)
			}
		})
	}
}

func TestQuadtreeInBoundMatching(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	q := New(geo.Bound{Max: geo.Point{5, 5}})
	pointers := []dataPointer{
		{geo.Point{0, 0}, false},
		{geo.Point{1, 1}, true},
		{geo.Point{2, 2}, false},
		{geo.Point{3, 3}, true},
		{geo.Point{4, 4}, false},
		{geo.Point{5, 5}, true},
	}
	for _, p := range pointers {
		if err := q.Add(p); err != nil {
			t.Fatalf("unexpected error for %v: %v", p, err)
		}
	}

	filters := map[bool]FilterFunc{
		false: nil,
		true:  func(p geo.Pointer) bool { return p.(dataPointer).visible },
	}

	cases := []struct {
		name     string
		filtered bool
		expected []geo.Point
	}{
		{
			name:     "unfiltered",
			filtered: false,
			expected: []geo.Point{{0, 0}, {1, 1}, {2, 2}},
		},
		{
			name:     "filtered",
			filtered: true,
			expected: []geo.Point{{1, 1}},
		},
	}

	var v []geo.Pointer
	bound := geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{2, 2}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v = q.InBoundMatching(v, bound, filters[tc.filtered])
			if len(v) != len(tc.expected) {
				t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
			}

			result := make([]geo.Point, 0)
			for _, p := range v {
				result = append(result, p.Point())
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i][0] < result[j][0]
			})

			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i][0] < tc.expected[j][0]
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Log(result)
				t.Log(tc.expected)
				t.Errorf("incorrect results")
			}
		})
	}
}

func TestQuadtreeInBound_Random(t *testing.T) {
	mp := geo.MultiPoint{}
	r := rand.New(rand.NewSource(43))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		mp = append(mp, geo.Point{r.Float64(), r.Float64()})
		if err := qt.Add(mp[i]); err != nil {
			t.Fatalf("unexpected error for %v: %v", mp[i], err)
		}
	}

	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		b := geo.Bound{Min: p, Max: p}
		b = b.Pad(0.1)
		ps := qt.InBound(nil, b)

		// find the right answer brute force
		var list []geo.Pointer
		for _, p := range mp {
			if b.Contains(p) {
				list = append(list, p)
			}
		}

		if len(list) != len(ps) {
			t.Errorf("index: %d, lengths not equal %v != %v", i, len(list), len(ps))
		}
	}
}

func TestQuadtreeAdd(t *testing.T) {
	p := geo.Point{}
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 10; i++ {
		// should be able to insert the same point over and over.
		if err := qt.Add(p); err != nil {
			t.Fatalf("unexpected error for %v: %v", p, err)
		}
	}
}

func countNodes(n *node) (c int) {
	if n != nil {
		c = 1 + countNodes(n.Children[0]) + countNodes(n.Children[1]) + countNodes(n.Children[2]) + countNodes(n.Children[3])
	}
	return
}
