package quadtree

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

type PExtra struct {
	p  geo.Point
	id string
}

func (p *PExtra) Point() geo.Point {
	return p.p
}

func (p *PExtra) String() string {
	return fmt.Sprintf("%v: %v", p.id, p.p)
}

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
		if err := qt.Add(p); err != nil {
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
		t.Fatalf("unexpected error: %e", err)
	}

	if err := qt.Add(dataPointer{geo.Point{1, 1}, true}); err != nil {
		t.Fatalf("unexpected error: %e", err)
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

func TestQuadtreeRemove(t *testing.T) {
	mp := geo.MultiPoint{}
	r := rand.New(rand.NewSource(42))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	for i := 0; i < 1000; i++ {
		mp = append(mp, geo.Point{r.Float64(), r.Float64()})
		if err := qt.Add(mp[i]); err != nil {
			t.Fatalf("unexpected error for %v: %v", mp[i], err)
		}
	}

	for i := 0; i < 1000; i += 3 {
		qt.Remove(mp[i], nil)
		mp[i] = geo.Point{-10000, -10000}
	}

	// make sure finding still works for 1000 random points
	for i := 0; i < 1000; i++ {
		p := geo.Point{r.Float64(), r.Float64()}
		f := qt.Find(p)
		_, j := planar.DistanceFromWithIndex(mp, p)
		if e := mp[j]; !e.Equal(f.Point()) {
			t.Errorf("index: %d, unexpected point %v != %v", i, e, f.Point())
		}
	}
}

func TestQuadtreeRemoveAndAdd_inOrder(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	r := rand.New(rand.NewSource(seed))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	p1 := &PExtra{p: geo.Point{r.Float64(), r.Float64()}, id: "1"}
	p2 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "2"}
	p3 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "3"}
	if err := qt.Add(p1); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	if err := qt.Add(p2); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	if err := qt.Add(p3); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	// rm 3
	found := qt.Remove(p3, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p3.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	// leaf node doesn't actually get removed
	if c := countNodes(qt.root); c != 3 {
		t.Errorf("incorrect number of nodes: %v != 3", c)
	}

	// 3 again
	found = qt.Remove(p3, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p3.id
	})
	if found {
		t.Errorf("should not find already removed node")
	}

	// rm 2
	found = qt.Remove(p2, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p2.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 2 {
		t.Errorf("incorrect number of nodes: %v != 2", c)
	}

	// rm 1
	found = qt.Remove(p1, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p1.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}
}

func TestQuadtreeRemoveAndAdd_sameLoc(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	r := rand.New(rand.NewSource(seed))
	qt := New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
	p1 := &PExtra{p: geo.Point{r.Float64(), r.Float64()}, id: "1"}
	p2 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "2"}
	p3 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "3"}
	p4 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "4"}
	p5 := &PExtra{p: geo.Point{p1.p[0], p1.p[1]}, id: "5"}

	if err := qt.Add(p1); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	if err := qt.Add(p2); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	if err := qt.Add(p3); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	// remove middle point
	found := qt.Remove(p2, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p2.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 2 {
		t.Errorf("incorrect number of nodes: %v != 2", c)
	}

	// remove first point
	found = qt.Remove(p1, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p1.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}

	// add a 4th point
	if err := qt.Add(p4); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	// remove third point
	found = qt.Remove(p3, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p3.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}

	// add a 5th point
	if err := qt.Add(p5); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	// remove the 5th point
	found = qt.Remove(p5, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p5.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	// 5 is a tail point, so its not does not actually get removed
	if c := countNodes(qt.root); c != 2 {
		t.Errorf("incorrect number of nodes: %v != 2", c)
	}

	// add a 3th point again
	if err := qt.Add(p3); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	// should reuse the tail point left by p5
	if c := countNodes(qt.root); c != 2 {
		t.Errorf("incorrect number of nodes: %v != 2", c)
	}

	// remove p4/root
	found = qt.Remove(p4, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p4.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}

	// remove p3/root
	found = qt.Remove(p3, func(p geo.Pointer) bool {
		return p.(*PExtra).id == p3.id
	})
	if !found {
		t.Error("didn't find/remove point")
	}

	// just the root, can't remove it
	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}

	// add back a point to be put in the root
	if err := qt.Add(p3); err != nil {
		t.Fatalf("unexpected error: %e", err)
	}

	if c := countNodes(qt.root); c != 1 {
		t.Errorf("incorrect number of nodes: %v != 1", c)
	}
}

func TestQuadtreeRemoveAndAdd_random(t *testing.T) {
	const perRun = 300
	const runs = 10
	var id int
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	r := rand.New(rand.NewSource(seed))
	bounds := geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{3000, 3000}}
	qt := New(bounds)
	points := make([]*PExtra, 0, 3000)
	for i := 0; i < runs; i++ {
		for j := 0; j < perRun; j++ {
			x := r.Int63n(30)
			y := r.Int63n(30)
			id++
			p := &PExtra{p: geo.Point{float64(x), float64(y)}, id: fmt.Sprintf("%d", id)}
			if err := qt.Add(p); err != nil {
				t.Fatalf("unexpected error for %v: %v", p, err)
			}

			points = append(points, p)

		}

		for j := 0; j < perRun/2; j++ {
			k := r.Int() % len(points)
			remP := points[k]
			points = append(points[:k], points[k+1:]...)
			qt.Remove(remP, func(p geo.Pointer) bool {
				return p.(*PExtra).id == remP.id
			})
		}
	}

	left := len(qt.InBound(nil, bounds))
	expected := runs * perRun / 2
	if left != expected {
		t.Errorf("incorrect number of points in tree: %d != %d", left, expected)
	}
}

func TestQuadtreeKNearest(t *testing.T) {
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
		point    geo.Point
		expected []geo.Point
	}{
		{
			name:     "unfiltered",
			filtered: false,
			point:    geo.Point{0.1, 0.1},
			expected: []geo.Point{{0, 0}, {1, 1}},
		},
		{
			name:     "filtered",
			filtered: true,
			point:    geo.Point{0.1, 0.1},
			expected: []geo.Point{{1, 1}, {3, 3}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.filtered {
				v := q.KNearest(nil, tc.point, 2)
				if len(v) != len(tc.expected) {
					t.Errorf("incorrect response length: %d != %d", len(v), len(tc.expected))
				}
			}

			v := q.KNearestMatching(nil, tc.point, 2, filters[tc.filtered])
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

func TestQuadtreeKNearest_sorted(t *testing.T) {
	q := New(geo.Bound{Max: geo.Point{5, 5}})
	for i := 0; i <= 5; i++ {
		if err := q.Add(geo.Point{float64(i), float64(i)}); err != nil {
			t.Fatalf("unexpected error: %e", err)
		}
	}

	nearest := q.KNearest(nil, geo.Point{2.25, 2.25}, 5)
	expected := []geo.Point{{2, 2}, {3, 3}, {1, 1}, {4, 4}, {0, 0}}
	for i, p := range expected {
		if n := nearest[i].Point(); !n.Equal(p) {
			t.Errorf("incorrect point %d: %v", i, n)
		}
	}
}

func TestQuadtreeKNearest_sorted2(t *testing.T) {
	q := New(geo.Bound{Max: geo.Point{8, 8}})
	for i := 0; i <= 7; i++ {
		if err := q.Add(geo.Point{float64(i), float64(i)}); err != nil {
			t.Fatalf("unexpected error: %e", err)
		}
	}

	nearest := q.KNearest(nil, geo.Point{5.25, 5.25}, 3)
	expected := []geo.Point{{5, 5}, {6, 6}, {4, 4}}
	for i, p := range expected {
		if n := nearest[i].Point(); !n.Equal(p) {
			t.Errorf("incorrect point %d: %v", i, n)
		}
	}
}

func TestQuadtreeKNearest_DistanceLimit(t *testing.T) {
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
		distance float64
		point    geo.Point
		expected []geo.Point
	}{
		{
			name:     "filtered",
			filtered: true,
			distance: 5,
			point:    geo.Point{0.1, 0.1},
			expected: []geo.Point{{1, 1}, {3, 3}},
		},
		{
			name:     "unfiltered",
			filtered: false,
			distance: 1,
			point:    geo.Point{0.1, 0.1},
			expected: []geo.Point{{0, 0}},
		},
	}

	var v []geo.Pointer
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v = q.KNearestMatching(v, tc.point, 5, filters[tc.filtered], tc.distance)
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

func countNodes(n *node) (c int) {
	if n != nil {
		c = 1 + countNodes(n.Children[0]) + countNodes(n.Children[1]) + countNodes(n.Children[2]) + countNodes(n.Children[3])
	}
	return
}
