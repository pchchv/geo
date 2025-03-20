package ewkb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/binary"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	_ sql.Scanner  = &GeometryScanner{}
	_ driver.Value = value{}
)

// GeometryScanner scans the results of sql queries.
// It can be used as a scan destination:
//
//	var s wkb.GeometryScanner
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(&s)
//	...
//	if s.Valid {
//	  // use s.Geometry
//	  // use s.SRID
//	} else {
//	  // NULL value
//	}
type GeometryScanner struct {
	sridInPrefix bool
	g            interface{}
	SRID         int
	Geometry     geo.Geometry
	Valid        bool // Valid is true if the geometry is not NULL
}

// Scanner returns a GeometryScanner that can scan sql query results.
// The geometryScanner.Geometry attribute will be set to the value.
// If g is non-nil, it MUST be a pointer to an geo.Geometry type like a Point or LineString.
// In that case the value will be written to g and the Geometry attribute.
//
//	var p geo.Point
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(wkb.Scanner(&p))
//	...
//	// use p
//
// If the value may be null check Valid first:
//
//	var point geo.Point
//	s := wkb.Scanner(&point)
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(s)
//	...
//	if s.Valid {
//	  // use p
//	} else {
//	  // NULL value
//	}
func Scanner(g interface{}) *GeometryScanner {
	return &GeometryScanner{g: g}
}

// ScannerPrefixSRID scans ewkb data were the SRID is in the first 4 bytes of the data.
// Databases like mysql/mariadb use this as their raw format.
// This method should only be used when working with such a database.
//
//	var p geo.Point
//	err := db.QueryRow("SELECT latlon FROM foo WHERE id=?", id).Scan(wkb.PrefixSRIDScanner(&p))
//
// However, it is recommended to covert to wkb explicitly using something like:
//
//	var srid int
//	var p geo.Point
//	err := db.QueryRow("SELECT ST_SRID(latlon), ST_AsBinary(latlon) FROM foo WHERE id=?", id).
//		Scan(&srid, wkb.Scanner(&p))
//
// https://dev.mysql.com/doc/refman/5.7/en/gis-data-formats.html
func ScannerPrefixSRID(g interface{}) *GeometryScanner {
	return &GeometryScanner{sridInPrefix: true, g: g}
}

// Scan scans the input []byte data into a geometry.
// This could be into the geo geometry type pointer or,
// if nil, the scanner.Geometry attribute.
func (s *GeometryScanner) Scan(d interface{}) error {
	var srid int
	var data interface{}
	s.Geometry = nil
	s.Valid = false
	data = d
	if s.sridInPrefix {
		if raw, ok := d.([]byte); !ok {
			return ErrUnsupportedDataType
		} else if raw == nil {
			return nil
		} else if len(raw) < 5 {
			return ErrNotEWKB
		} else {
			srid = int(binary.LittleEndian.Uint32(raw))
			data = raw[4:]
		}
	}

	g, embeddedSRID, valid, err := wkbcommon.Scan(s.g, data)
	if err != nil {
		return mapCommonError(err)
	} else if embeddedSRID != 0 {
		srid = embeddedSRID
	}

	s.Geometry = g
	s.SRID = srid
	s.Valid = valid

	return nil
}

type value struct {
	srid int
	v    geo.Geometry
}

// Value creates a driver.Valuer that will EWKB the geometry into the database query.
// db.Exec("INSERT INTO table (point_column) VALUES (?)", ewkb.Value(p, 4326))
func Value(g geo.Geometry, srid int) driver.Valuer {
	return value{srid: srid, v: g}
}

func (v value) Value() (driver.Value, error) {
	if val, err := Marshal(v.v, v.srid); val == nil {
		return nil, err
	} else {
		return val, err
	}
}

type valuePrefixSRID struct {
	srid int
	v    geo.Geometry
}

// ValuePrefixSRID creates a driver.Valuer that
// will WKB the geometry but add the srid as a 4 byte prefix.
// db.Exec("INSERT INTO table (point_column) VALUES (?)", ewkb.Value(p, 4326))
func ValuePrefixSRID(g geo.Geometry, srid int) driver.Valuer {
	return valuePrefixSRID{srid: srid, v: g}
}

func (v valuePrefixSRID) Value() (driver.Value, error) {
	if val, err := Marshal(v.v, 0); val == nil || err != nil {
		return nil, err
	} else {
		data := make([]byte, 4, 4+len(val))
		binary.LittleEndian.PutUint32(data, uint32(v.srid))
		return append(data, val...), nil
	}
}
