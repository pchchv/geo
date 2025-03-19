package wkb

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	_ sql.Scanner  = &GeometryScanner{}
	_ driver.Value = value{}
)

// GeometryScanner scans the results of sql queries,
// it can be used as a scan destination:
//
//	s := &wkb.GeometryScanner{}
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(s)
//	...
//	if s.Valid {
//	  // use s.Geometry
//	} else {
//	  // NULL value
//	}
type GeometryScanner struct {
	g        interface{}
	Geometry geo.Geometry
	Valid    bool // Valid is true if the geometry is not NULL
}

// Scan scanes the input []byte data into a geometry.
// This could be into the geo geometry type pointer or, if nil,
// the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	if d == nil {
		return nil
	}

	data, ok := d.([]byte)
	if !ok {
		return ErrUnsupportedDataType
	}

	s.Geometry = nil
	s.Valid = false
	if g, _, valid, err := wkbcommon.Scan(s.g, d); err == wkbcommon.ErrNotWKBHeader {
		var e error
		g, _, valid, e = wkbcommon.Scan(s.g, data[4:])
		if e != wkbcommon.ErrNotWKBHeader {
			err = e // nil or incorrect type, e.g. decoding line string
		}
	} else if err != nil {
		return mapCommonError(err)
	} else {
		s.Geometry = g
		s.Valid = valid
	}

	return nil
}

type value struct {
	v geo.Geometry
}

// Value creates a driver.Valuer that will WKB the geometry into the database query.
func Value(g geo.Geometry) driver.Valuer {
	return value{v: g}

}
