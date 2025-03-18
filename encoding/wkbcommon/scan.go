package wkbcommon

import (
	"errors"

	"github.com/pchchv/geo"
)

var (
	ErrNotWKB              = errors.New("wkbcommon: invalid data")         // returned when unmarshalling WKB and the data is not valid
	ErrNotWKBHeader        = errors.New("wkbcommon: invalid header data")  // returned when unmarshalling first few bytes and there is an issue
	ErrIncorrectGeometry   = errors.New("wkbcommon: incorrect geometry")   // returned when unmarshalling WKB data into the wrong type (e. g. linestring into a point)
	ErrUnsupportedGeometry = errors.New("wkbcommon: unsupported geometry") // returned when geometry type is not supported by this package
)

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
