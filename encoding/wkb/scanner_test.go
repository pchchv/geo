package wkb

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/pchchv/geo"
)

var SRID = []byte{215, 15, 0, 0}

func TestScanNil(t *testing.T) {
	s := Scanner(nil)
	if err := s.Scan(testPointData); err != nil {
		t.Fatalf("scan error: %e", err)
	}

	if !geo.Equal(s.Geometry, testPoint) {
		t.Errorf("incorrect geometry: %v != %v", s.Geometry, testPoint)
	}

	t.Run("scan nil data", func(t *testing.T) {
		var p geo.Point
		s := Scanner(&p)
		if err := s.Scan(nil); err != nil {
			t.Errorf("should noop for nil data: %e", err)
		}

		if s.Valid {
			t.Errorf("valid should be false for nil values")
		}
	})

	t.Run("scan nil byte interface", func(t *testing.T) {
		var b []byte
		var p geo.Point
		s := Scanner(&p)
		if err := s.Scan(b); err != nil {
			t.Errorf("should noop for nil data: %e", err)
		}

		if s.Valid {
			t.Errorf("valid should be false for nil values")
		}
	})

	t.Run("unknown geometry type", func(t *testing.T) {
		s := Scanner(nil)
		b := []byte{
			//01    02    03    04    05    06    07    08
			0x01, 0x08, 0x00, 0x00, 0x00, // CircularString type
			0x46, 0x81, 0xF6, 0x23, 0x2E, 0x4A, 0x5D, 0xC0,
			0x03, 0x46, 0x1B, 0x3C, 0xAF, 0x5B, 0x40, 0x40,
		}
		if err := s.Scan(b); err != ErrUnsupportedGeometry {
			t.Errorf("incorrect error: %e != %e", err, ErrUnsupportedGeometry)
		}

		if s.Valid {
			t.Errorf("valid should be false errors")
		}
	})
}

func TestScanHexData(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Point
	}{
		{
			name:     "point lower case",
			data:     []byte(`\x0101000000e0d57267266e4840b22ac24d46b50240`),
			expected: geo.Point{48.860547, 2.338513},
		},
		{
			name:     "point upper case",
			data:     []byte(`\x0101000000e0d57267266e4840b22ac24d46b50240`),
			expected: geo.Point{48.860547, 2.338513},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Point
			s := Scanner(&p)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(p)
				t.Log(tc.expected)
			}
		})
	}
}

func TestScanHexData_errors(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Point
	}{
		{
			name:     "not hex data",
			data:     []byte(`\xZZ0101000000e0d57267266e4840b22ac24d46b50240`),
			expected: geo.Point{48.860547, 2.338513},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Point
			s := Scanner(&p)
			if err := s.Scan(tc.data); err == nil {
				t.Fatalf("should have error, but no error")
			}
		})
	}
}

func TestScanPoint(t *testing.T) {
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
			name:     "point with MySQL SRID",
			data:     append(SRID, testPointData...),
			expected: testPoint,
		},
		{
			name:     "single multi-point",
			data:     testMultiPointSingleData,
			expected: testMultiPointSingle[0],
		},
		{
			name:     "single multi-point with MySQL SRID",
			data:     append(SRID, testMultiPointSingleData...),
			expected: testMultiPointSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Point
			s := Scanner(&p)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(p)
				t.Log(tc.expected)
			}

			if p != s.Geometry {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanPoint_errors(t *testing.T) {
	// error conditions
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb, too short",
			data: []byte{0, 0, 0, 0, 1, 192, 94, 157, 24, 227, 60, 152, 15, 64, 66, 222, 128, 39},
			err:  ErrNotWKB,
		},
		{
			name: "invalid first byte",
			data: []byte{3, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64},
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect geometry",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect geometry with MySQL SRID",
			data: append(SRID, testLineStringData...),
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Point
			s := Scanner(&p)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanMultiPoint(t *testing.T) {
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
			name:     "multi point with MySQL SRID",
			data:     append(SRID, testMultiPointData...),
			expected: testMultiPoint,
		},
		{
			name:     "point should covert to multi point",
			data:     testPointData,
			expected: geo.MultiPoint{testPoint},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mp geo.MultiPoint
			s := Scanner(&mp)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !mp.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mp)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(mp, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanMultiPoint_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like line string",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 1, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mp geo.MultiPoint
			s := Scanner(&mp)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanLineString(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.LineString
	}{
		{
			name:     "line string",
			data:     testLineStringData,
			expected: testLineString,
		},
		{
			name:     "line string with MySQL SRID",
			data:     append(SRID, testLineStringData...),
			expected: testLineString,
		},
		{
			name:     "single multi line string",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ls geo.LineString
			s := Scanner(&ls)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !ls.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(ls)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(ls, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanLineString_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like multi point",
			data: testMultiPointData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 2, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ls geo.LineString
			s := Scanner(&ls)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanMultiLineString(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.MultiLineString
	}{
		{
			name:     "line string",
			data:     testLineStringData,
			expected: geo.MultiLineString{testLineString},
		},
		{
			name:     "multi line string",
			data:     testMultiLineStringData,
			expected: testMultiLineString,
		},
		{
			name:     "multi line string with MySQL SRID",
			data:     append(SRID, testMultiLineStringData...),
			expected: testMultiLineString,
		},
		{
			name:     "single multi line string",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mls geo.MultiLineString
			s := Scanner(&mls)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !mls.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mls)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(mls, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanMultiLineString_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like multi point",
			data: testMultiPointData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 5, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mls geo.MultiLineString
			s := Scanner(&mls)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanRing(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Ring
	}{
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: testPolygon[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var r geo.Ring
			s := Scanner(&r)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !r.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(r)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(r, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanRing_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 1, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var r geo.Ring
			s := Scanner(&r)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanPolygon(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Polygon
	}{
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: testPolygon,
		},
		{
			name:     "polygon with MySQL SRID",
			data:     append(SRID, testPolygonData...),
			expected: testPolygon,
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle[0],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Polygon
			s := Scanner(&p)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !p.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(p)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(p, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanPolygon_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 3, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p geo.Polygon
			s := Scanner(&p)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanMultiPolygon(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.MultiPolygon
	}{
		{
			name:     "multi polygon",
			data:     testMultiPolygonData,
			expected: testMultiPolygon,
		},
		{
			name:     "multi polygon with MySQL SRID",
			data:     append(SRID, testMultiPolygonData...),
			expected: testMultiPolygon,
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle,
		},
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: geo.MultiPolygon{testPolygon},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mp geo.MultiPolygon
			s := Scanner(&mp)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !mp.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(mp)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(mp, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanMultiPolygon_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 6, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var mp geo.MultiPolygon
			s := Scanner(&mp)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanCollection(t *testing.T) {
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var c geo.Collection
			s := Scanner(&c)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !c.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(c)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(c, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanCollection_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "does not like line strings",
			data: testLineStringData,
			err:  ErrIncorrectGeometry,
		},
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 7, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var c geo.Collection
			s := Scanner(&c)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestScanBound(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected geo.Bound
	}{
		{
			name:     "point",
			data:     testPointData,
			expected: testPoint.Bound(),
		},
		{
			name:     "multi point",
			data:     testMultiPointData,
			expected: testMultiPoint.Bound(),
		},
		{
			name:     "single multi point",
			data:     testMultiPointSingleData,
			expected: testMultiPointSingle.Bound(),
		},
		{
			name:     "linestring",
			data:     testLineStringData,
			expected: testLineString.Bound(),
		},
		{
			name:     "multi linestring",
			data:     testMultiLineStringData,
			expected: testMultiLineString.Bound(),
		},
		{
			name:     "single multi linestring",
			data:     testMultiLineStringSingleData,
			expected: testMultiLineStringSingle.Bound(),
		},
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: testPolygon.Bound(),
		},
		{
			name:     "multi polygon",
			data:     testMultiPolygonData,
			expected: testMultiPolygon.Bound(),
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle.Bound(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b geo.Bound
			s := Scanner(&b)
			if err := s.Scan(tc.data); err != nil {
				t.Fatalf("scan error: %e", err)
			}

			if !b.Equal(tc.expected) {
				t.Errorf("unequal data")
				t.Log(b)
				t.Log(tc.expected)
			}

			if !reflect.DeepEqual(b, s.Geometry) {
				t.Errorf("should set to scanner's geometry attribute")
			}

			if !s.Valid {
				t.Errorf("should set valid to true")
			}
		})
	}
}

func TestScanBound_errors(t *testing.T) {
	cases := []struct {
		name string
		data interface{}
		err  error
	}{
		{
			name: "incorrect data",
			data: 123,
			err:  ErrUnsupportedDataType,
		},
		{
			name: "not wkb",
			data: []byte{0, 0, 0, 0, 1, 192, 94},
			err:  ErrNotWKB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b geo.Bound
			s := Scanner(&b)
			if err := s.Scan(tc.data); err != tc.err {
				t.Errorf("incorrect error: %v != %v", err, tc.err)
			}

			if s.Geometry != nil {
				t.Errorf("geometry should be nil on errors")
			}

			if s.Valid {
				t.Errorf("valid should be false on errors")
			}
		})
	}
}

func TestValue(t *testing.T) {
	t.Run("marshalls geometry", func(t *testing.T) {
		val, err := Value(testPoint).Value()
		if err != nil {
			t.Errorf("value error: %e", err)
		}

		if !bytes.Equal(val.([]byte), testPointData) {
			t.Errorf("incorrect marshal")
			t.Log(val)
			t.Log(testPointData)
		}
	})

	t.Run("nil value in should set nil value", func(t *testing.T) {
		val, err := Value(nil).Value()
		if err != nil {
			t.Errorf("value error: %e", err)
		}

		if val != nil {
			t.Errorf("should be nil value: %[1]T, %[1]v", val)
		}
	})
}

func TestValue_nil(t *testing.T) {
	var (
		mp    geo.MultiPoint
		ls    geo.LineString
		mls   geo.MultiLineString
		r     geo.Ring
		poly  geo.Polygon
		mpoly geo.MultiPolygon
		c     geo.Collection
	)

	cases := []struct {
		name string
		geom geo.Geometry
	}{
		{
			name: "nil multi point",
			geom: mp,
		},
		{
			name: "nil line string",
			geom: ls,
		},
		{
			name: "nil multi line string",
			geom: mls,
		},
		{
			name: "nil ring",
			geom: r,
		},
		{
			name: "nil polygon",
			geom: poly,
		},
		{
			name: "nil multi polygon",
			geom: mpoly,
		},
		{
			name: "nil collection",
			geom: c,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := Value(tc.geom).Value()
			if err != nil {
				t.Errorf("value error: %e", err)
			}

			if val != nil {
				t.Errorf("should be nil value: %[1]T, %[1]v", val)
			}
		})
	}
}
