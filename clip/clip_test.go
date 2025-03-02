package clip

import (
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

func TestInternalLine(t *testing.T) {
	cases := []struct {
		name   string
		bound  geo.Bound
		input  geo.LineString
		output geo.MultiLineString
	}{
		{
			name:  "clip line",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{30, 30}},
			input: geo.LineString{
				{-10, 10}, {10, 10}, {10, -10}, {20, -10}, {20, 10}, {40, 10},
				{40, 20}, {20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20}, {-10, 20},
			},
			output: geo.MultiLineString{
				{{0, 10}, {10, 10}, {10, 0}},
				{{20, 0}, {20, 10}, {30, 10}},
				{{30, 20}, {20, 20}, {20, 30}},
				{{10, 30}, {10, 20}, {5, 20}, {0, 20}},
			},
		},
		{
			name:  "clips line crossign many times",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{20, 20}},
			input: geo.LineString{
				{10, -10}, {10, 30}, {20, 30}, {20, -10},
			},
			output: geo.MultiLineString{
				{{10, 0}, {10, 20}},
				{{20, 20}, {20, 0}},
			},
		},
		{
			name:  "no changes if all inside",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{20, 20}},
			input: geo.LineString{
				{1, 1}, {2, 2}, {3, 3},
			},
			output: geo.MultiLineString{
				{{1, 1}, {2, 2}, {3, 3}},
			},
		},
		{
			name:  "empty if nothing in bound",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{2, 2}},
			input: geo.LineString{
				{10, 10}, {20, 20}, {30, 30},
			},
			output: nil,
		},
		{
			name:  "floating point example",
			bound: geo.Bound{Min: geo.Point{-91.93359375, 42.29356419217009}, Max: geo.Point{-91.7578125, 42.42345651793831}},
			input: geo.LineString{
				{-86.66015624999999, 42.22851735620852}, {-81.474609375, 38.51378825951165},
				{-85.517578125, 37.125286284966776}, {-85.8251953125, 38.95940879245423},
				{-90.087890625, 39.53793974517628}, {-91.93359375, 42.32606244456202},
				{-86.66015624999999, 42.22851735620852},
			},
			output: geo.MultiLineString{
				{
					{-91.91208030440808, 42.29356419217009},
					{-91.93359375, 42.32606244456202},
					{-91.7578125, 42.3228109416169},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := line(tc.bound, tc.input, false)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}

func TestInternalRing(t *testing.T) {
	cases := []struct {
		name   string
		bound  geo.Bound
		input  geo.Ring
		output geo.Ring
	}{
		{
			name:  "clips polygon",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{30, 30}},
			input: geo.Ring{
				{-10, 10}, {0, 10}, {10, 10}, {10, 5}, {10, -5},
				{10, -10}, {20, -10}, {20, 10}, {40, 10}, {40, 20},
				{20, 20}, {20, 40}, {10, 40}, {10, 20}, {5, 20},
				{-10, 20}},
			// note: we allow duplicate points if polygon endpoints are
			// on the box boundary.
			output: geo.Ring{
				{0, 10}, {0, 10}, {10, 10}, {10, 5}, {10, 0},
				{20, 0}, {20, 10}, {30, 10}, {30, 20}, {20, 20},
				{20, 30}, {10, 30}, {10, 20}, {5, 20}, {0, 20},
			},
		},
		{
			name:  "completely inside bound",
			bound: geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{10, 10}},
			input: geo.Ring{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
			output: geo.Ring{
				{3, 3}, {5, 3}, {5, 5}, {3, 5}, {3, 3},
			},
		},
		{
			name:  "completely around bound",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{2, 2}},
			input: geo.Ring{
				{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0},
			},
			output: geo.Ring{{1, 2}, {1, 1}, {2, 1}, {2, 2}, {1, 2}},
		},
		{
			name:  "completely around touching corners",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{3, 3}},
			input: geo.Ring{
				{0, 2}, {2, 0}, {4, 2}, {2, 4}, {0, 2},
			},
			output: geo.Ring{{1, 1}, {1, 1}, {3, 1}, {3, 1}, {3, 3}, {3, 3}, {1, 3}, {1, 3}, {1, 1}},
		},
		{
			name:  "around but cut corners",
			bound: geo.Bound{Min: geo.Point{0.5, 0.5}, Max: geo.Point{3.5, 3.5}},
			input: geo.Ring{
				{0, 2}, {2, 4}, {4, 2}, {2, 0}, {0, 2},
			},
			output: geo.Ring{{0.5, 2.5}, {1.5, 3.5}, {2.5, 3.5}, {3.5, 2.5}, {3.5, 1.5}, {2.5, 0.5}, {1.5, 0.5}, {0.5, 1.5}, {0.5, 2.5}},
		},
		{
			name:  "unclosed ring",
			bound: geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{4, 4}},
			input: geo.Ring{
				{2, 0}, {3, 0}, {3, 5}, {2, 5},
			},
			output: geo.Ring{{3, 1}, {3, 4}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ring(tc.bound, tc.input)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect clip")
				t.Logf("%v", result)
				t.Logf("%v", tc.output)
			}
		})
	}
}
