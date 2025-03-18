package wkbcommon

import (
	"math"

	"github.com/pchchv/geo"
)

func (e *Encoder) writePolygon(p geo.Polygon, srid int) (err error) {
	if err = e.writeTypePrefix(polygonType, len(p), srid); err != nil {
		return
	}

	for _, r := range p {
		e.order.PutUint32(e.buf, uint32(len(r)))
		if _, err = e.w.Write(e.buf[:4]); err != nil {
			return
		}

		for _, p := range r {
			e.order.PutUint64(e.buf, math.Float64bits(p[0]))
			e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
			if _, err = e.w.Write(e.buf); err != nil {
				return
			}
		}
	}

	return nil
}

func (e *Encoder) writeMultiPolygon(mp geo.MultiPolygon, srid int) (err error) {
	if err = e.writeTypePrefix(multiPolygonType, len(mp), srid); err != nil {
		return
	}

	for _, p := range mp {
		if err = e.Encode(p, 0); err != nil {
			return
		}
	}

	return nil
}
