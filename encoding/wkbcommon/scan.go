package wkbcommon

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/pchchv/geo"
)

var (
	ErrNotWKB              = errors.New("wkbcommon: invalid data")              // returned when unmarshalling WKB and the data is not valid
	ErrNotWKBHeader        = errors.New("wkbcommon: invalid header data")       // returned when unmarshalling first few bytes and there is an issue
	ErrIncorrectGeometry   = errors.New("wkbcommon: incorrect geometry")        // returned when unmarshalling WKB data into the wrong type (e. g. linestring into a point)
	ErrUnsupportedGeometry = errors.New("wkbcommon: unsupported geometry")      // returned when geometry type is not supported by this package
	ErrUnsupportedDataType = errors.New("wkbcommon: scan value must be []byte") // returned when the data type is not []byte
)

// Scan scans the input []byte data into a geometry.
// This can be a geo geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func Scan(g, d interface{}) (geo.Geometry, int, bool, error) {
	if d == nil {
		return nil, 0, false, nil
	}

	data, ok := d.([]byte)
	if !ok {
		return nil, 0, false, ErrUnsupportedDataType
	} else if data == nil {
		return nil, 0, false, nil
	} else if len(data) < 5 {
		return nil, 0, false, ErrNotWKB
	}

	// go-pg will return ST_AsBinary(*) data as `\xhexencoded` which
	// needs to be converted to true binary for further decoding
	// code detects the \x prefix and then converts the rest from Hex to binary
	if data[0] == byte('\\') && data[1] == byte('x') {
		n, err := hex.Decode(data, data[2:])
		if err != nil {
			return nil, 0, false, fmt.Errorf("thought the data was hex with prefix, but it is not: %v", err)
		}
		data = data[:n]
	}

	// also possible is just straight hex encoded
	// in this case the bo bit can be '0x00' or '0x01'
	if data[0] == '0' && (data[1] == '0' || data[1] == '1') {
		n, err := hex.Decode(data, data)
		if err != nil {
			return nil, 0, false, fmt.Errorf("thought the data was hex, but it is not: %v", err)
		}
		data = data[:n]
	}

	switch g := g.(type) {
	case nil:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return nil, 0, false, err
		}

		return m, srid, true, nil
	case *geo.Point:
		p, srid, err := ScanPoint(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = p
		return p, srid, true, nil
	case *geo.MultiPoint:
		m, srid, err := ScanMultiPoint(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = m
		return m, srid, true, nil
	case *geo.LineString:
		l, srid, err := ScanLineString(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = l
		return l, srid, true, nil
	case *geo.MultiLineString:
		m, srid, err := ScanMultiLineString(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = m
		return m, srid, true, nil
	case *geo.Ring:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return nil, 0, false, err
		}

		if p, ok := m.(geo.Polygon); ok && len(p) == 1 {
			*g = p[0]
			return p[0], srid, true, nil
		}

		return nil, 0, false, ErrIncorrectGeometry
	case *geo.Polygon:
		p, srid, err := ScanPolygon(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = p
		return p, srid, true, nil
	case *geo.MultiPolygon:
		m, srid, err := ScanMultiPolygon(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = m
		return m, srid, true, nil
	case *geo.Collection:
		c, srid, err := ScanCollection(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = c
		return c, srid, true, nil
	case *geo.Bound:
		m, srid, err := Unmarshal(data)
		if err != nil {
			return nil, 0, false, err
		}

		*g = m.Bound()
		return *g, srid, true, nil
	default:
		return nil, 0, false, ErrIncorrectGeometry
	}
}

// ScanPoint takes binary wkb and decodes it into a point.
func ScanPoint(data []byte) (geo.Point, int, error) {
	order, typ, srid, geomData, err := unmarshalByteOrderType(data)
	if err != nil {
		return geo.Point{}, 0, err
	}

	switch typ {
	case pointType:
		if p, err := unmarshalPoint(order, geomData); err != nil {
			return geo.Point{}, 0, err
		} else {
			return p, srid, nil
		}
	case multiPointType:
		if mp, err := unmarshalMultiPoint(order, geomData); err != nil {
			return geo.Point{}, 0, err
		} else if len(mp) == 1 {
			return mp[0], srid, nil
		}
	}

	return geo.Point{}, 0, ErrIncorrectGeometry
}

// ScanMultiPoint takes binary wkb and decodes it into a multi-point.
func ScanMultiPoint(data []byte) (geo.MultiPoint, int, error) {
	m, srid, err := Unmarshal(data)
	if err != nil {
		return nil, 0, err
	}

	switch p := m.(type) {
	case geo.Point:
		return geo.MultiPoint{p}, srid, nil
	case geo.MultiPoint:
		return p, srid, nil
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanLineString takes binary wkb and decodes it into a line string.
func ScanLineString(data []byte) (geo.LineString, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case lineStringType:
		if ls, err := unmarshalLineString(order, data); err != nil {
			return nil, 0, err
		} else {
			return ls, srid, nil
		}
	case multiLineStringType:
		if mls, err := unmarshalMultiLineString(order, data); err != nil {
			return nil, 0, err
		} else if len(mls) == 1 {
			return mls[0], srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanMultiLineString takes binary wkb and decodes it into a multi-line string.
func ScanMultiLineString(data []byte) (geo.MultiLineString, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case lineStringType:
		if ls, err := unmarshalLineString(order, data); err != nil {
			return nil, 0, err
		} else {
			return geo.MultiLineString{ls}, srid, nil
		}
	case multiLineStringType:
		if ls, err := unmarshalMultiLineString(order, data); err != nil {
			return nil, 0, err
		} else {
			return ls, srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanPolygon takes binary wkb and decodes it into a polygon.
func ScanPolygon(data []byte) (geo.Polygon, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case polygonType:
		if p, err := unmarshalPolygon(order, data); err != nil {
			return nil, 0, err
		} else {
			return p, srid, nil
		}
	case multiPolygonType:
		if mp, err := unmarshalMultiPolygon(order, data); err != nil {
			return nil, 0, err
		} else if len(mp) == 1 {
			return mp[0], srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanMultiPolygon takes binary wkb and decodes it into a multi-polygon.
func ScanMultiPolygon(data []byte) (geo.MultiPolygon, int, error) {
	order, typ, srid, data, err := unmarshalByteOrderType(data)
	if err != nil {
		return nil, 0, err
	}

	switch typ {
	case polygonType:
		if p, err := unmarshalPolygon(order, data); err != nil {
			return nil, 0, err
		} else {
			return geo.MultiPolygon{p}, srid, nil
		}
	case multiPolygonType:
		if mp, err := unmarshalMultiPolygon(order, data); err != nil {
			return nil, 0, err
		} else {
			return mp, srid, nil
		}
	}

	return nil, 0, ErrIncorrectGeometry
}

// ScanCollection takes binary wkb and decodes it into a collection.
func ScanCollection(data []byte) (geo.Collection, int, error) {
	m, srid, err := NewDecoder(bytes.NewReader(data)).Decode()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return nil, 0, ErrNotWKB
	} else if err != nil {
		return nil, 0, err
	}

	switch p := m.(type) {
	case geo.Collection:
		return p, srid, nil
	default:
		return nil, 0, ErrIncorrectGeometry
	}
}
