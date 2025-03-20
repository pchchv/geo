package ewkb

import "github.com/pchchv/geo"

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

// ScannerPrefixSRID will scan ewkb data were the SRID is in the first 4 bytes of the data.
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
