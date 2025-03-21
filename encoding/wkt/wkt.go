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
