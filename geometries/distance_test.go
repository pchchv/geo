package geometries

import (
	"math"
	"testing"

	"github.com/pchchv/geo"
)

const epsilon = 1e-6

func TestDistance(t *testing.T) {
	p1 := geo.Point{-1.8444, 53.1506}
	p2 := geo.Point{0.1406, 52.2047}
	if d := Distance(p1, p2); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = geo.Point{0.5, 30}
	p2 = geo.Point{-0.5, 30}
	dFast := Distance(p1, p2)
	p1 = geo.Point{179.5, 30}
	p2 = geo.Point{-179.5, 30}
	if d := Distance(p1, p2); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestDistanceHaversine(t *testing.T) {
	p1 := geo.Point{-1.8444, 53.1506}
	p2 := geo.Point{0.1406, 52.2047}
	if d := DistanceHaversine(p1, p2); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}

	p1 = geo.Point{0.5, 30}
	p2 = geo.Point{-0.5, 30}
	dHav := DistanceHaversine(p1, p2)
	p1 = geo.Point{179.5, 30}
	p2 = geo.Point{-179.5, 30}
	if d := DistanceHaversine(p1, p2); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect distance, got %v", d)
	}
}

func TestMidpoint(t *testing.T) {
	answer := geo.Point{-0.841153, 52.68179432}
	m := Midpoint(geo.Point{-1.8444, 53.1506}, geo.Point{0.1406, 52.2047})
	if d := Distance(m, answer); d > 1 {
		t.Errorf("expected %v, got %v", answer, m)
	}
}

func TestPointAtBearingAndDistance(t *testing.T) {
	cases := []struct {
		name     string
		point    geo.Point
		bearing  float64
		distance float64
		expected geo.Point
	}{
		{
			name:     "simple",
			point:    geo.Point{-1.8444, 53.1506},
			bearing:  127.373,
			distance: 85194.89,
			expected: geo.Point{-0.841153, 52.68179432},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := PointAtBearingAndDistance(tc.point, tc.bearing, tc.distance)
			if d := DistanceHaversine(actual, tc.expected); d > 1 {
				t.Errorf("expected %v, got %v (%vm away)", tc.expected, actual, d)
			}
		})
	}

	t.Run("midpoint", func(t *testing.T) {
		a := geo.Point{-1.8444, 53.1506}
		b := geo.Point{0.1406, 52.2047}
		bearing := Bearing(a, b)
		distance := DistanceHaversine(a, b)
		p1 := PointAtBearingAndDistance(a, bearing, distance/2)
		p2 := Midpoint(a, b)
		if d := DistanceHaversine(p1, p2); d > epsilon {
			t.Errorf("expected %v to be within %vm of %v", p1, epsilon, p2)
		}
	})
}

func TestPointAtDistanceAlongLineWithSinglePoint(t *testing.T) {
	cases := []struct {
		name            string
		line            geo.LineString
		distance        float64
		expectedPoint   geo.Point
		expectedBearing float64
	}{
		{
			name: "with single point",
			line: geo.LineString{
				geo.Point{-1.8444, 53.1506},
			},
			distance:        9000,
			expectedPoint:   geo.Point{-1.8444, 53.1506},
			expectedBearing: 0,
		},
		{
			name: "with minimal points",
			line: geo.LineString{
				geo.Point{-1.8444, 53.1506},
				geo.Point{0.1406, 52.2047},
			},
			distance:      85194.89,
			expectedPoint: geo.Point{-0.841153, 52.68179432},
			expectedBearing: Bearing(
				geo.Point{-1.8444, 53.1506},
				geo.Point{0.1406, 52.2047},
			),
		},
		{
			name: "with single point",
			line: geo.LineString{
				geo.Point{-1.8444, 53.1506},
				geo.Point{-0.8411, 52.6817},
				geo.Point{0.1406, 52.2047},
			},
			distance:      90000,
			expectedPoint: geo.Point{-0.78526, 52.65506},
			expectedBearing: Bearing(
				geo.Point{-0.8411, 52.6817},
				geo.Point{0.1406, 52.2047},
			),
		},
		{
			name: "past end of line",
			line: geo.LineString{
				geo.Point{-1.8444, 53.1506},
				geo.Point{-0.8411, 52.6817},
				geo.Point{0.1406, 52.2047},
			},
			distance:      200000,
			expectedPoint: geo.Point{0.1406, 52.2047},
			expectedBearing: Bearing(
				geo.Point{-0.8411, 52.6817},
				geo.Point{0.1406, 52.2047},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualPoint, actualBearing := PointAtDistanceAlongLine(tc.line, tc.distance)
			if d := DistanceHaversine(actualPoint, tc.expectedPoint); d > 1 {
				t.Errorf("point %v != %v", actualPoint, tc.expectedPoint)
			}

			if d := math.Abs(actualBearing - tc.expectedBearing); d > 1 {
				t.Errorf("bearing %v != %v", tc.expectedBearing, actualBearing)
			}
		})
	}
}

func TestPointAtDistanceAlongLineWithEmptyLineString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("PointAtDistanceAlongLine did not panic")
		}
	}()

	line := geo.LineString{}
	PointAtDistanceAlongLine(line, 90000)
}

func TestBearing(t *testing.T) {
	p1 := geo.Point{0, 0}
	p2 := geo.Point{0, 1}
	if d := Bearing(p1, p2); d != 0 {
		t.Errorf("expected 0, got %f", d)
	}

	if d := Bearing(p2, p1); d != 180 {
		t.Errorf("expected 180, got %f", d)
	}

	p1 = geo.Point{0, 0}
	p2 = geo.Point{1, 0}
	if d := Bearing(p1, p2); d != 90 {
		t.Errorf("expected 90, got %f", d)
	}

	if d := Bearing(p2, p1); d != -90 {
		t.Errorf("expected -90, got %f", d)
	}

	p1 = geo.Point{-1.8444, 53.1506}
	p2 = geo.Point{0.1406, 52.2047}
	if d := Bearing(p1, p2); math.Abs(127.373351-d) > epsilon {
		t.Errorf("point, bearingTo got %f", d)
	}
}
