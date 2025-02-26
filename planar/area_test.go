package planar

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestCentroidArea(t *testing.T) {
	for _, g := range geo.AllGeometries {
		CentroidArea(g)
	}
}

func TestCentroidArea_MultiPoint(t *testing.T) {
	mp := geo.MultiPoint{{0, 0}, {1, 1.5}, {2, 0}}
	centroid, area := CentroidArea(mp)
	expected := geo.Point{1, 0.5}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	}

	if area != 0 {
		t.Errorf("area should be 0: %f", area)
	}
}

func TestCentroidArea_LineString(t *testing.T) {
	cases := []struct {
		name   string
		ls     geo.LineString
		result geo.Point
	}{
		{
			name:   "simple",
			ls:     geo.LineString{{0, 0}, {3, 4}},
			result: geo.Point{1.5, 2},
		},
		{
			name:   "empty line",
			ls:     geo.LineString{},
			result: geo.Point{0, 0},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ls); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestCentroidArea_MultiLineString(t *testing.T) {
	cases := []struct {
		name   string
		ls     geo.MultiLineString
		result geo.Point
	}{
		{
			name:   "simple",
			ls:     geo.MultiLineString{{{0, 0}, {3, 4}}},
			result: geo.Point{1.5, 2},
		},
		{
			name:   "two lines",
			ls:     geo.MultiLineString{{{0, 0}, {0, 1}}, {{1, 0}, {1, 1}}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "multiple empty lines",
			ls:     geo.MultiLineString{{{1, 0}}, {{2, 1}}},
			result: geo.Point{1.5, 0.5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ls); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestCentroid_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   geo.Ring
		result geo.Point
	}{
		{
			name:   "triangle, cw",
			ring:   geo.Ring{{0, 0}, {1, 3}, {2, 0}, {0, 0}},
			result: geo.Point{1, 1},
		},
		{
			name:   "triangle, ccw",
			ring:   geo.Ring{{0, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: geo.Point{1, 1},
		},
		{
			name:   "square, cw",
			ring:   geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "non-closed square, cw",
			ring:   geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "triangle, ccw",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: geo.Point{0.5, 0.5},
		},
		{
			name:   "redudent points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: geo.Point{1, 1},
		},
		{
			name: "3 points",
			ring: geo.Ring{{0, 0}, {1, 0}, {0, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e8
				tc.ring[i][1] -= 1e8
			}

			tc.result[0] += 1e8
			tc.result[1] -= 1e8
			if c, _ := CentroidArea(tc.ring); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestArea_Ring(t *testing.T) {
	cases := []struct {
		name   string
		ring   geo.Ring
		result float64
	}{
		{
			name:   "simple box, ccw",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			ring:   geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: -1,
		},
		{
			name:   "even number of points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0.4, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "3 points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {0, 0}},
			result: 0.0,
		},
		{
			name:   "4 points",
			ring:   geo.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			result: 0.5,
		},
		{
			name:   "6 points",
			ring:   geo.Ring{{1, 1}, {2, 1}, {2, 1.5}, {2, 2}, {1, 2}, {1, 1}},
			result: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, val := CentroidArea(tc.ring); val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that is recenters to deal with roundoff
			for i := range tc.ring {
				tc.ring[i][0] += 1e15
				tc.ring[i][1] -= 1e15
			}

			if _, val := CentroidArea(tc.ring); val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}

			// check that are rendant last point is implicit
			tc.ring = tc.ring[:len(tc.ring)-1]
			if _, val := CentroidArea(tc.ring); val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

// +-+ +-+
// | | | |
// | +-+ |
// |     |
// +-----+
func TestCentroid_RingAdv(t *testing.T) {
	ring := geo.Ring{{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5}, {2, 1}, {3, 1}, {3, 0}, {0, 0}}
	centroid, area := CentroidArea(ring)
	expected := geo.Point{1.5, 0.45}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	} else if area != -2.5 {
		t.Errorf("incorrect area: %v != 2.5", area)
	}
}

func TestCentroidArea_Polygon(t *testing.T) {
	t.Run("polygon with hole", func(t *testing.T) {
		r1 := geo.Ring{{0, 0}, {4, 0}, {4, 3}, {0, 3}, {0, 0}}
		r1.Reverse()
		r2 := geo.Ring{{2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 1}}
		poly := geo.Polygon{r1, r2}
		centroid, area := CentroidArea(poly)
		if !centroid.Equal(geo.Point{21.5 / 11.0, 1.5}) {
			t.Errorf("%v", 21.5/11.0)
			t.Errorf("incorrect centroid: %v", centroid)
		} else if area != 11 {
			t.Errorf("incorrect area: %v != 11", area)
		}
	})

	t.Run("collapsed", func(t *testing.T) {
		e := geo.Point{0.5, 1}
		if c, _ := CentroidArea(geo.Polygon{{{0, 1}, {1, 1}, {0, 1}}}); !c.Equal(e) {
			t.Errorf("incorrect point: %v != %v", c, e)
		}
	})

	t.Run("empty right half", func(t *testing.T) {
		poly := geo.Polygon{
			{{0, 0}, {4, 0}, {4, 4}, {0, 4}, {0, 0}},
			{{2, 0}, {2, 4}, {4, 4}, {4, 0}, {2, 0}},
		}

		centroid, area := CentroidArea(poly)
		if v := (geo.Point{1, 2}); !centroid.Equal(v) {
			t.Errorf("incorrect centroid: %v != %v", centroid, v)
		} else if area != 8 {
			t.Errorf("incorrect area: %v != 8", area)
		}
	})
}

func TestCentroidArea_Bound(t *testing.T) {
	b := geo.Bound{Min: geo.Point{0, 2}, Max: geo.Point{1, 3}}
	centroid, area := CentroidArea(b)
	expected := geo.Point{0.5, 2.5}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	} else if area != 1 {
		t.Errorf("incorrect area: %f != 1", area)
	}

	b = geo.Bound{Min: geo.Point{0, 2}, Max: geo.Point{0, 2}}
	centroid, area = CentroidArea(b)
	expected = geo.Point{0, 2}
	if !centroid.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", centroid, expected)
	} else if area != 0 {
		t.Errorf("area should be zero: %f", area)
	}
}
