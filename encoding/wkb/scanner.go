package wkb

import "github.com/pchchv/geo"

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
