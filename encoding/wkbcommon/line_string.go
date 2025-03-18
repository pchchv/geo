package wkbcommon

import (
	"math"

	"github.com/pchchv/geo"
)

func (e *Encoder) writeLineString(ls geo.LineString, srid int) (err error) {
	if err = e.writeTypePrefix(lineStringType, len(ls), srid); err != nil {
		return
	}

	for _, p := range ls {
		e.order.PutUint64(e.buf, math.Float64bits(p[0]))
		e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
		if _, err = e.w.Write(e.buf); err != nil {
			return
		}
	}

	return nil
}

func (e *Encoder) writeMultiLineString(mls geo.MultiLineString, srid int) (err error) {
	if err = e.writeTypePrefix(multiLineStringType, len(mls), srid); err != nil {
		return
	}

	for _, ls := range mls {
		if err = e.Encode(ls, 0); err != nil {
			return
		}
	}

	return nil
}
