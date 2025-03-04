package simplify

import (
	"testing"

	"github.com/pchchv/geo"
)

func TestDoubleTriangleArea(t *testing.T) {
	expected := 30.0
	ls := geo.LineString{{2, 5}, {5, 1}, {-4, 3}}
	cases := []struct {
		name       string
		i1, i2, i3 int
	}{
		{
			name: "order 1",
			i1:   0, i2: 1, i3: 2,
		},
		{
			name: "order 2",
			i1:   0, i2: 2, i3: 1,
		},
		{
			name: "order 3",
			i1:   1, i2: 2, i3: 0,
		},
		{
			name: "order 4",
			i1:   1, i2: 0, i3: 2,
		},
		{
			name: "order 5",
			i1:   2, i2: 0, i3: 1,
		},
		{
			name: "order 6",
			i1:   2, i2: 1, i3: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			area := doubleTriangleArea(ls, tc.i1, tc.i2, tc.i3)
			if area != expected {
				t.Errorf("incorrect area: %v != %v", area, expected)
			}
		})
	}
}
