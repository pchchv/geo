package ewkb

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
		0x01, 0x07, 0x00, 0x00, 0x20,
		0xE6, 0x10, 0x00, 0x00,
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
		srid     int
		data     []byte
		expected geo.Collection
	}{
		{
			name:     "large",
			srid:     123,
			data:     MustMarshal(large, 123),
			expected: large,
		},
		{
			name:     "collection with point",
			data:     MustDecodeHex("0107000020e6100000010000000101000000000000000000f03f0000000000000040"),
			srid:     4326,
			expected: geo.Collection{geo.Point{1, 2}},
		},
		{
			name: "collection with point and line",
			data: MustDecodeHex("0020000007000010e60000000200000000013ff000000000000040000000000000000000000002000000023ff0000000000000400000000000000040080000000000004010000000000000"),
			srid: 4326,
			expected: geo.Collection{
				geo.Point{1, 2},
				geo.LineString{{1, 2}, {3, 4}},
			},
		},
		{
			name: "collection with point and line and polygon",
			data: MustDecodeHex("0107000020e6100000030000000101000000000000000000f03f0000000000000040010200000002000000000000000000f03f00000000000000400000000000000840000000000000104001030000000300000004000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001840000000000000f03f000000000000004004000000000000000000264000000000000028400000000000002a400000000000002c400000000000002e4000000000000030400000000000002640000000000000284004000000000000000000354000000000000036400000000000003740000000000000384000000000000039400000000000003a4000000000000035400000000000003640"),
			srid: 4326,
			expected: geo.Collection{
				geo.Point{1, 2},
				geo.LineString{{1, 2}, {3, 4}},
				geo.Polygon{
					{{1, 2}, {3, 4}, {5, 6}, {1, 2}},
					{{11, 12}, {13, 14}, {15, 16}, {11, 12}},
					{{21, 22}, {23, 24}, {25, 26}, {21, 22}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.srid, tc.data)
		})
	}
}
