package wkb

import (
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	testPoint     = geo.Point{-117.15906619141342, 32.71628524142945}
	testPointData = []byte{
		//01    02    03    04    05    06    07    08
		0x01, 0x01, 0x00, 0x00, 0x00,
		0x46, 0x81, 0xF6, 0x23, 0x2E, 0x4A, 0x5D, 0xC0,
		0x03, 0x46, 0x1B, 0x3C, 0xAF, 0x5B, 0x40, 0x40,
	}
	testMultiPoint     = geo.MultiPoint{{10, 40}, {40, 30}, {20, 20}, {30, 10}}
	testMultiPointData = []byte{
		0x01, 0x04, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00, // Number of Points (4)
		0x01,                   // Byte Order Little
		0x01, 0x00, 0x00, 0x00, // Type Point (1)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
		0x01,                   // Byte Order Little
		0x01, 0x00, 0x00, 0x00, // Type Point (1)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 30
		0x01,                   // Byte Order Little
		0x01, 0x00, 0x00, 0x00, // Type Point (1)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
		0x01,                   // Byte Order Little
		0x01, 0x00, 0x00, 0x00, // Type Point (1)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y4 10
	}
	testMultiPointSingle     = geo.MultiPoint{{10, 40}}
	testMultiPointSingleData = []byte{
		0x01, 0x04, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, // Number of Points (4)
		0x01,                   // Byte Order Little
		0x01, 0x00, 0x00, 0x00, // Type Point (1)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
	}
)

func TestPoint(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Point
	}{
		{
			name:     "point",
			data:     testPointData,
			expected: testPoint,
		},
		{
			name:     "little endian",
			data:     []byte{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			expected: geo.Point{-122.4546440212, 37.7382859071},
		},
		{
			name:     "big endian",
			data:     []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39, 17, 11, 205},
			expected: geo.Point{-122.4546440212, 37.7382859071},
		},
		{
			name:     "another point",
			data:     []byte{1, 1, 0, 0, 0, 253, 104, 56, 101, 110, 114, 87, 192, 192, 9, 133, 8, 56, 50, 64, 64},
			expected: geo.Point{-93.787988, 32.392335},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}

func TestPointToHex(t *testing.T) {
	cases := []struct {
		name     string
		data     geo.Point
		expected string
	}{
		{
			name:     "point",
			data:     geo.Point{1, 2},
			expected: "0101000000000000000000f03f0000000000000040",
		},
		{
			name:     "zero point",
			data:     geo.Point{0, 0},
			expected: "010100000000000000000000000000000000000000",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := MustMarshalToHex(tc.data)
			if s != tc.expected {
				t.Errorf("incorrect hex: %v", s)
			}
		})
	}
}

func TestMultiPoint(t *testing.T) {
	large := geo.MultiPoint{}
	for i := 0; i < wkbcommon.MaxPointsAlloc+100; i++ {
		large = append(large, geo.Point{float64(i), float64(-i)})
	}

	cases := []struct {
		name     string
		data     []byte
		expected geo.MultiPoint
	}{
		{
			name:     "multi point",
			data:     testMultiPointData,
			expected: testMultiPoint,
		},
		{
			name:     "single multi point",
			data:     testMultiPointSingleData,
			expected: testMultiPointSingle,
		},
		{
			name:     "large multi point",
			data:     MustMarshal(large),
			expected: large,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}
