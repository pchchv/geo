package wkt

import (
	"bytes"
	"fmt"

	"github.com/pchchv/geo"
)

func writeLineString(buf *bytes.Buffer, ls geo.LineString) {
	buf.WriteByte('(')
	for i, p := range ls {
		if i != 0 {
			buf.WriteByte(',')
		}

		fmt.Fprintf(buf, "%g %g", p[0], p[1])
	}
	buf.WriteByte(')')
}

func wkt(buf *bytes.Buffer, geom geo.Geometry) {
	switch g := geom.(type) {
	case geo.Point:
		fmt.Fprintf(buf, "POINT(%g %g)", g[0], g[1])
	case geo.MultiPoint:
		if len(g) == 0 {
			buf.Write([]byte(`MULTIPOINT EMPTY`))
			return
		}

		buf.Write([]byte(`MULTIPOINT(`))
		for i, p := range g {
			if i != 0 {
				buf.WriteByte(',')
			}

			fmt.Fprintf(buf, "(%g %g)", p[0], p[1])
		}

		buf.WriteByte(')')
	case geo.LineString:
		if len(g) == 0 {
			buf.Write([]byte(`LINESTRING EMPTY`))
			return
		}

		buf.Write([]byte(`LINESTRING`))
		writeLineString(buf, g)
	case geo.MultiLineString:
		if len(g) == 0 {
			buf.Write([]byte(`MULTILINESTRING EMPTY`))
			return
		}

		buf.Write([]byte(`MULTILINESTRING(`))
		for i, ls := range g {
			if i != 0 {
				buf.WriteByte(',')
			}

			writeLineString(buf, ls)
		}

		buf.WriteByte(')')
	case geo.Ring:
		wkt(buf, geo.Polygon{g})
	case geo.Polygon:
		if len(g) == 0 {
			buf.Write([]byte(`POLYGON EMPTY`))
			return
		}

		buf.Write([]byte(`POLYGON(`))
		for i, r := range g {
			if i != 0 {
				buf.WriteByte(',')
			}
			writeLineString(buf, geo.LineString(r))
		}

		buf.WriteByte(')')
	case geo.MultiPolygon:
		if len(g) == 0 {
			buf.Write([]byte(`MULTIPOLYGON EMPTY`))
			return
		}

		buf.Write([]byte(`MULTIPOLYGON(`))
		for i, p := range g {
			if i != 0 {
				buf.WriteByte(',')
			}

			buf.WriteByte('(')
			for j, r := range p {
				if j != 0 {
					buf.WriteByte(',')
				}

				writeLineString(buf, geo.LineString(r))
			}

			buf.WriteByte(')')
		}

		buf.WriteByte(')')
	case geo.Collection:
		if len(g) == 0 {
			buf.Write([]byte(`GEOMETRYCOLLECTION EMPTY`))
			return
		}

		buf.Write([]byte(`GEOMETRYCOLLECTION(`))
		for i, c := range g {
			if i != 0 {
				buf.WriteByte(',')
			}

			wkt(buf, c)
		}

		buf.WriteByte(')')
	case geo.Bound:
		wkt(buf, g.ToPolygon())
	default:
		panic("unsupported type")
	}
}
