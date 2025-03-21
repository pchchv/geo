package wkb

import (
	"testing"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	testCollection = geo.Collection{
		geo.Point{4, 6},
		geo.LineString{{4, 6}, {7, 10}},
	}
	testCollectionData = []byte{
		//01    02    03    04    05    06    07    08
		0x01, 0x07, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, // Number of Geometries in Collection
		0x01,                   // Byte order marker little
		0x01, 0x00, 0x00, 0x00, // Type (1) Point
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40, // X1 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40, // Y1 6
		0x01,                   // Byte order marker little
		0x02, 0x00, 0x00, 0x00, // Type (2) Line
		0x02, 0x00, 0x00, 0x00, // Number of Points (2)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40, // X1 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40, // Y1 6
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c, 0x40, // X2 7
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y2 10
	}
)

func TestCollection(t *testing.T) {
	large := geo.Collection{}
	for i := 0; i < wkbcommon.MaxMultiAlloc+100; i++ {
		large = append(large, geo.Point{float64(i), float64(-i)})
	}

	cases := []struct {
		name     string
		data     []byte
		expected geo.Collection
	}{
		{
			name:     "collection",
			data:     testCollectionData,
			expected: testCollection,
		},
		{
			name:     "large",
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
